package LinkServer

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"net"
	"os"
	"time"
)

func init() {
	connMap = make(map[*net.Conn]*UserConnection)
}

var connMap map[*net.Conn]*UserConnection

const (
	HEADER_SIGN = 0x5F5F
)

var etcdClient *clientv3.Client

type UserConnection struct {
	UserId           int64
	Session          string
	RawIp            int64
	ProxyId          int64 //haproxy
	HallConnection   net.Conn
	GameConnection   net.Conn
	ClientConnection net.Conn
}

func (uc *UserConnection) Close() {
	_, f := connMap[&uc.ClientConnection]
	if f {
		delete(connMap, &uc.ClientConnection)
	}
	uc.ClientConnection.Close()
	uc.UserId = 0
	if uc.HallConnection != nil {
		uc.HallConnection.Close()
		uc.HallConnection = nil
	}
	if uc.GameConnection != nil {
		uc.GameConnection.Close()
		uc.GameConnection = nil
	}
	uc.RawIp = 0
	uc.ProxyId = 0
	uc.Session = ""
}

func StartProxy(port int32) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	f, err := os.Create(fmt.Sprintf("./logs/proxy-%d.log", port))
	if err != nil {
		fmt.Println("Error when creating log file:", err)
		os.Exit(1)
	}

	log.SetOutput(f)
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Println("Error Listening:", err)
		os.Exit(1)
	}
	resolveLocalIp()
	startEtcd()
	getServers()
	watchServers(l.Addr().String())

	defer func() {
		l.Close()
		etcdClient.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			os.Exit(1)
		}
		log.Printf("Receivied message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		initSession(conn)
		go handleRequest(conn)
	}
}

func initSession(conn net.Conn) {
	session := genSessionStr(16)
	userConnection := &UserConnection{
		Session:          session,
		ClientConnection: conn,
	}
	connMap[&conn] = userConnection
}

func genSessionStr(len int) string {
	key := make([]byte, len)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("Generate session key error:%s", err)
		return ""
	}
	return fmt.Sprintf("%x", key)
}

func handleRequest(conn net.Conn) {
	uc, _ := connMap[&conn]
	defer func() {
		uc.Close()
	}()
	b := make([]byte, 1024)
	n, err := conn.Read(b)
	if err != nil {
		log.Printf("remote client %s data read error %s", conn.RemoteAddr(), err)
		return
	}

	if n == 0 || n < binary.Size(Header{}) {
		log.Printf("invalid stream from client: %s,size:%d", conn.RemoteAddr(), n)
		return
	}

	var header Header
	var buf *bytes.Buffer
	buf.ReadFrom(conn)
	header.FromBuf(buf)
	if header.Sign != HEADER_SIGN {
		log.Printf("client:%s invalid header sign:%d", conn.RemoteAddr(), header.Sign)
		return
	}
	if header.Ver != 1 {
		fmt.Printf("client:%s invalid header sign:%d", conn.RemoteAddr(), header.Sign)
		return
	}
	//todo check crc
	//checkData := buf.Bytes()[4:]
	//sum := md5.Sum(checkData)
	//if sum != header.Crc {
	//
	//}
	//todo decrypto
	_, f := connMap[&conn]
	if !f {
		fmt.Println("Can't find UserConnection for client:", conn.RemoteAddr())
	}
	switch header.MainId {
	case MSG_CLIENT_TO_PROXY:
		break
	case MSG_CLIENT_TO_GAME:
	case MSG_CLIENT_TO_HALL:
		interHeader := InternalPreHeader{}
		data := bytes.Buffer{}
		data.Write(interHeader.ToBuf().Bytes())
		data.Write(buf.Bytes())
		if header.MainId == MSG_CLIENT_TO_HALL {
			sendToHall(uc, data.Bytes())
		} else {
			sendToGame(uc, data.Bytes())
		}
		break
	}
}

func sendToHall(conn *UserConnection, data []byte) {
	if conn.HallConnection == nil {

	}
}

func sendToGame(conn *UserConnection, data []byte) {

}

func handleGetAesKey() {

}

func startEtcd() {
	var err error
	etcdClient, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:3379", "localhost:4379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln("connect failed, err:", err)
		return
	}
	log.Println("connect successed")
}

func getServers() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	resp, err := etcdClient.Get(ctx, "servers/proxyservers/")
	cancel()
	if err != nil {
		log.Println("get failed, err:", err)
		return
	}
	for _, ev := range resp.Kvs {
		log.Printf("%s : %s \n", ev.Key, ev.Value)
	}
}

func watchServers(addr string) {
	etcdClient.Put(context.Background(), "servers/proxyservers/", addr)
	for {
		rch := etcdClient.Watch(context.Background(), "servers/proxyservers/")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				log.Printf("Event[type:%s key:%q value:%q] \n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}

func resolveLocalIp() {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalln("Get net interface error:", err)
		return
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Println("Loop net interface error:", err)
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP

			}
			log.Println("Get local ip address:", ip.String())
		}
	}
}
