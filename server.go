// @Title  main.go
// @Description A highly interactive honeypot supporting redis protocol
// @Author  Cy 2021.04.08
package main

import (
	"bytes"
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"github.com/walu/resp"
	"log"
)

var (
	err error
)

type RedisServer struct {
	server *gev.Server
}

func NewRedisServer(address string, proto string, loopsnum int) (server *RedisServer, err error) {
	Serv := new(RedisServer)
	Serv.server, err = gev.NewServer(Serv,
		gev.Address(address),
		gev.Network(proto),
		gev.NumLoops(loopsnum))
	if err != nil {
		return nil, err
		panic(err)
	}
	return Serv, nil
}

func (s *RedisServer) Start() {
	s.server.Start()
}

func (s *RedisServer) Stop() {
	s.server.Stop()
}

func (s *RedisServer) OnConnect(c *connection.Connection) {
	log.Println(" OnConnect ： ", c.PeerAddr())
}

func (s *RedisServer) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	out = data
	command := bytes.NewReader(data)
	cmd, err := resp.ReadCommand(command)
	if err != nil {
		out = data
	}
	switch cmd.Name() {
	case "ping":
		out = []byte("+PONG\r\n")
		return
	case "info":
		return
	case "":
		return
	default:
		out = []byte("-ERR wrong number of arguments for " + cmd.Name() + " command\r\n")
		return
	}
	return
}

func (s *RedisServer) OnClose(c *connection.Connection) {
	log.Println(c.PeerAddr(), "Closed")
}
