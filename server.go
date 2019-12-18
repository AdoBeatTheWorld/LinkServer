package linkserver

import (
	"encoding/json"
	"fmt"
	"github.com/funny/link"
	"github.com/funny/link/codec"
	"gitlab.com/adoontheway/HttpServer/db"
	"gitlab.com/adoontheway/HttpServer/redis"
	"gitlab.com/adoontheway/LinkProto"
	"io/ioutil"
	"log"
	"os"
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
	Start(config.Port)
}

func Start(port int32) {
	json := codec.Json()
	json.Register(proto.AddReq{})
	json.Register(proto.AddRsp{})
	server, err := link.Listen("tcp", fmt.Sprintf(":%d", port), json, 0, link.HandlerFunc(serverSessionLoop))
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
