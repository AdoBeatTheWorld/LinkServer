package LinkServer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type UserConnection struct {
	UserId int64
	Session string
	RawIp int64
	ProxyId int64//haproxy
	HallConnection IConnection
	GameConnection IConnection
}

type IConnection struct {
	Conn *net.Conn
}

func StartProxy()  {
	l,err := net.Listen("tcp","0.0.0.0:8888")
	if err != nil {
		fmt.Println("Error Listening:",err)
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
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn)  {
	defer conn.Close()
	b := make([]byte,1024)
	n,err := conn.Read(b)
	if err != nil {
		fmt.Printf("remote client %s data read error %s",conn.RemoteAddr(), err)
		return
	}

	if n == 0 || n < binary.Size(Header{}) {
		fmt.Printf("invalid stream from client: %s,size:%d", conn.RemoteAddr(), n)
	}
	var header Header
	var buf *bytes.Buffer
	buf.ReadFrom(conn)
	header.FromBuf(buf)
	//todo make sum
	//todo check crc
	//todo decrypto
	switch header.MainId {
	case MSG_CLIENT_TO_PROXY:
		break;
	case MSG_CLIENT_TO_GAME:
		break;
	case MSG_CLIENT_TO_HALL:
		break;
	}
}

func onUserLogin(conn *IConnection, session string)  {

}