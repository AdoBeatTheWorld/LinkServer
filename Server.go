package LinkServer

import (
	"LinkProto"
	"encoding/json"
	"gitlab.com/adoontheway/HttpServer/db"
	"gitlab.com/adoontheway/HttpServer/redis"
	"io/ioutil"
	"log"
	"os"

	"github.com/funny/link"
	"github.com/funny/link/codec"
)

type ServerConfig struct {
	Port      int32  `json:"port"`
	Rocketmq  int32  `json:"rocketmq"`
	Consul    string `json:"consul"`
	LogLevel  string `json:"log_level"`
	LogPath   string `json:"log_path"`
	DB        string `json:"db"`
	Redis     string `json:"redis"`
	RedisPass string `json:"redis_pass"`
}

func InitServerFromConfig(configPath string) {
	config, err := ReadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
	db.InitMongoConnector(config.DB)
	redis.InitRedisPool(config.Redis, config.RedisPass)
	Start()
}

func Start() {
	json := codec.Json()
	json.Register(LinkProto.AddReq{})
	json.Register(LinkProto.AddRsp{})

	server, err := link.Listen("tcp", "0.0.0.0:8888", json, 0, link.HandlerFunc(serverSessionLoop))
	checkErr(err)
	server.Serve()
}

func serverSessionLoop(session *link.Session) {
	for {
		req, err := session.Receive()
		checkErr(err)

		err = session.Send(&LinkProto.AddRsp{
			req.(*LinkProto.AddReq).A + req.(*LinkProto.AddReq).B,
		})
		checkErr(err)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ReadConfig(filepath string) (*ServerConfig, error) {
	configfile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
		//utils.Zapper.Error(err.Error())
		return nil, err
	}
	defer configfile.Close()

	byteValue, err := ioutil.ReadAll(configfile)
	if err != nil {
		log.Fatal(err)
		//utils.Zapper.Error(err.Error())
		return nil, err
	}
	var config ServerConfig
	json.Unmarshal(byteValue, &config)
	return &config, nil
}
