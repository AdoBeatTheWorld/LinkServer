# LinkServer
这是一个golang游戏服务器的基础搭建项目  
先搭建一个简单的可用的服务器，然后逐步分解成一个分布式的服务器  
后续添加功能和支持如下：  
## etcd
分布式服务发现
## mongodb & mysql & redis
数据库支持
## rocketmq
消息队列支持
## haproxy
这个是运维的，用于负载均衡
## protobuf
通讯数据格式
## aes & rsa
通讯过程中需要用到的加密算法，或者其他算法亦可