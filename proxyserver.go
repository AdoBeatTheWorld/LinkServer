package LinkServer

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func init() {
	connMap = make(map[*net.Conn]*UserConnection)
}

var connMap map[*net.Conn]*UserConnection

const (
	HEADER_SIGN = 0x5F5F
)

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

func StartProxy() {
	l, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Println("Error Listening:", err)
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err)
			os.Exit(1)
		}
		fmt.Printf("Receivied message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
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
		fmt.Printf("Generate session key error:%s", err)
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
		fmt.Printf("remote client %s data read error %s", conn.RemoteAddr(), err)
		return
	}

	if n == 0 || n < binary.Size(Header{}) {
		fmt.Printf("invalid stream from client: %s,size:%d", conn.RemoteAddr(), n)
		return
	}

	var header Header
	var buf *bytes.Buffer
	buf.ReadFrom(conn)
	header.FromBuf(buf)
	if header.Sign != HEADER_SIGN {
		fmt.Printf("client:% invalid header sign:%d", conn.RemoteAddr(), header.Sign)
		return
	}
	if header.Ver != 1 {
		fmt.Printf("client:% invalid header sign:%d", conn.RemoteAddr(), header.Sign)
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
		interHeader := InternalPreHeader{}
		data := bytes.Buffer{}
		data.Write(interHeader.ToBuf().Bytes())
		data.Write(buf.Bytes())
		sendToGame(uc, data.Bytes())
		break
	case MSG_CLIENT_TO_HALL:
		interHeader := InternalPreHeader{}
		data := bytes.Buffer{}
		data.Write(interHeader.ToBuf().Bytes())
		data.Write(buf.Bytes())
		sendToHall(uc, data.Bytes())
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
