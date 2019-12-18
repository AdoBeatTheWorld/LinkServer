package linkserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type Service struct {
	Name    string
	Info    ServiceInfo
	stop    chan error
	leaseid clientv3.LeaseID
	client  *clientv3.Client
}

type ServiceInfo struct {
	IP string
}

func (s *Service) Start() error {
	ch, err := s.keepAlive()
	if err != nil {
		log.Fatal(err)
		return err
	}
	for {
		select {
		case err := <-s.stop:
			s.revoke()
			return err
		case <-s.client.Ctx().Done():
			return errors.New("service closed")
		case ka, ok := <-ch:
			if !ok {
				log.Println("keep alive channel closed")
				s.revoke()
				return nil
			} else {
				log.Printf("Recv reply from service:%s, ttl:%d", s.Name, ka.TTL)
			}

		}
	}
}

func (s *Service) Stop() {
	s.stop <- nil
}

func (s *Service) keepAlive() (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	info := &s.Info
	key := "services/" + s.Name
	value, _ := json.Marshal(info)
	resp, err := s.client.Grant(context.TODO(), 5)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	_, err = s.client.Put(context.TODO(), key, string(value), clientv3.WithLease(resp.ID))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	s.leaseid = resp.ID
	return s.client.KeepAlive(context.TODO(), resp.ID)
}

func (s *Service) revoke() error {
	_, err := s.client.Revoke(context.TODO(), s.leaseid)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Service:%s stoped\n", s.Name)
	return err
}

func NewService(name string, info ServiceInfo, endpoints []string) (*Service, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &Service{
		Name:   name,
		Info:   info,
		stop:   make(chan error),
		client: cli,
	}, err
}

type Master struct {
	Path   string
	Client *clientv3.Client
	Nodes  map[string]*Node
}

type Node struct {
	State bool
	Key   string
	Info  ServiceInfo
}

func NewMaster(endpoints []string, watchPath string) (*Master, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	master := &Master{
		Path:   watchPath,
		Nodes:  make(map[string]*Node),
		Client: cli,
	}
	go master.WatchNodes()
	return master, err
}

func (m *Master) AddNode(key string, info *ServiceInfo) {
	node := &Node{
		State: true,
		Key:   key,
		Info:  *info,
	}
	m.Nodes[node.Key] = node
}

func (m *Master) WatchNodes() {
	rch := m.Client.Watch(context.Background(), m.Path, clientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case clientv3.EventTypePut:
				fmt.Printf("[%s] %q : %q \n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				info := GetServiceInfo(ev)
				m.AddNode(string(ev.Kv.Key), info)
			case clientv3.EventTypeDelete:
				fmt.Printf("[%s] %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				delete(m.Nodes, string(ev.Kv.Key))
			}
		}
	}
}

func GetServiceInfo(ev *clientv3.Event) *ServiceInfo {
	info := &ServiceInfo{}
	err := json.Unmarshal([]byte(ev.Kv.Value), info)
	if err != nil {
		log.Println(err)
	}
	return info
}
