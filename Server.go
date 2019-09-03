package LinkServer

import (
	"log"
	"proto"

	"github.com/funny/link"
	"github.com/funny/link/codec"
)


func Start()  {
		json := codec.Json()
		json.Register(proto.AddReq{})
		json.Register(proto.AddRsp{})

		server, err := link.Listen("tcp", "0.0.0.0:8888", json, 0 , link.HandlerFunc(serverSessionLoop))
		checkErr(err)
		server.Serve()
}

func serverSessionLoop(session *link.Session) {
	for {
		req, err := session.Receive()
		checkErr(err)

		err = session.Send(&proto.AddRsp{
			req.(*proto.AddReq).A + req.(*proto.AddReq).B,
		})
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}