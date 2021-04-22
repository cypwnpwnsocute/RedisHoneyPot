// @Title  server.go
// @Description High Interaction Honeypot Solution for Redis protocol
// @Author  Cy 2021.04.08
package main

import (
	"bytes"
	"fmt"
	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/sirupsen/logrus"
	"github.com/walu/resp"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
)

type RedisServer struct {
	server  *gev.Server
	hashmap *hashmap.Map
	Config  *ini.File
	log     *logrus.Logger
}

func NewRedisServer(address string, proto string, loopsnum int) (server *RedisServer, err error) {
	Serv := new(RedisServer)
	Serv.hashmap = hashmap.New()
	config, err := LoadConfig("redis.conf")
	Serv.log = logrus.New()
	Serv.log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	if err != nil {
		panic(err)
	}
	Serv.Config = config
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
	s.log.WithFields(logrus.Fields{
		"action": "NewConnect",
		"addr":   c.PeerAddr(),
	}).Println()
}

func (s *RedisServer) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out []byte) {
	command := bytes.NewReader(data)
	if command.Len() == 2 {
		return
	}
	cmd, err := resp.ReadCommand(command)
	if err != nil {
		return
	}

	com := strings.ToLower(cmd.Name())

	s.log.WithFields(logrus.Fields{
		"action": strings.Join(cmd.Args, " "),
		"addr":   c.PeerAddr(),
	}).Println()

	switch com {
	case "ping":
		out = []byte("+PONG\r\n")
	case "info":
		info := ""
		for _, key := range s.Config.Section("info").KeyStrings() {
			info += fmt.Sprintf("%s:%s\r\n", key, s.Config.Section("info").Key(key))
		}
		out = []byte("$" + strconv.Itoa(len(info)) + "\r\n" + info + "\r\n")
	case "set":
		if len(cmd.Args) < 3 {
			out = []byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n")
		} else {
			s.hashmap.Put(cmd.Args[1], cmd.Args[2])
			out = []byte("+OK\r\n")
		}
	case "get":
		if len(cmd.Args) != 2 {
			out = []byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n")
		} else {
			v, bool := s.hashmap.Get(cmd.Args[1])
			if bool == true {
				out = []byte("+" + v.(string) + "\r\n")
			} else {
				out = []byte("+(nil)\r\n")
			}
		}
	case "del":
		if len(cmd.Args) < 2 {
			out = []byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n")
		} else {
			s.hashmap.Remove(cmd.Args[1])
			out = []byte("+(integer) 1\r\n")
		}
	case "exists":
		if len(cmd.Args) < 2 {
			out = []byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n")
		} else {
			_, bool := s.hashmap.Get(cmd.Args[1])
			if bool == true {
				out = []byte("+(integer) 1\r\n")
			} else {
				out = []byte("+(integer) 0\r\n")
			}
		}
	case "keys":
		if len(cmd.Args) != 2 {
			out = []byte("-ERR wrong number of arguments for '" + cmd.Args[0] + "' command\r\n")
		} else {
			if cmd.Args[1] == "*" {
				str := "*" + strconv.Itoa(s.hashmap.Size()) + "\r\n"
				for _, v := range s.hashmap.Keys() {
					str += "$" + strconv.Itoa(len(v.(string))) + "\r\n" + v.(string) + "\r\n"
				}
				out = []byte(str)
			} else {
				_, bool := s.hashmap.Get(cmd.Args[1])
				if bool == true {
					l := strconv.Itoa(len(cmd.Args[1]))
					out = []byte("*1\r\n$" + l + "\r\n" + cmd.Args[1] + "\r\n")
				} else {
					out = []byte("+(empty array)\r\n")
				}
			}
		}
	case "flushall":
		out = []byte("+OK\r\n")
	case "flushdb":
		out = []byte("+OK\r\n")
	case "save":
		out = []byte("+OK\r\n")
	case "select":
		out = []byte("+OK\r\n")
	case "dbsize":
		l := strconv.Itoa(s.hashmap.Size())
		out = []byte("+(integer) " + l + "\r\n")
	case "config":
		if cmd.Args[1] == "get" && len(cmd.Args) > 2 {
			if cmd.Args[2] != "*" {
				content := s.Config.Section("info").Key(cmd.Args[2]).String()
				if content == "" {
					out = []byte("+(empty array)\r\n")
				} else {
					l1 := strconv.Itoa(len(cmd.Args[2]))
					l2 := strconv.Itoa(len(content))
					out = []byte("*2\r\n$" + l1 + "\r\n" + cmd.Args[2] + "\r\n$" + l2 + "\r\n" + content + "\r\n")
				}
			} else {
				output := "*" + strconv.Itoa(len(s.Config.Section("info").KeyStrings())*2) + "\r\n"
				for _, key := range s.Config.Section("info").KeyStrings() {
					value := s.Config.Section("info").Key(key).String()
					output += "$" + strconv.Itoa(len(key)) + "\r\n" + key + "\r\n" + "$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n"
				}
				out = []byte(output)
			}
		} else if cmd.Args[1] == "set" && len(cmd.Args) > 2 {
			s.Config.Section("info").NewKey(cmd.Args[2], cmd.Args[3])
			out = []byte("+OK\r\n")
		} else {
			out = []byte("-ERR Unknown subcommand or wrong number of arguments for 'get'. Try CONFIG HELP.\r\n")
		}
	case "slaveof":
		if len(cmd.Args) < 3 {
			out = []byte("-ERR wrong number of arguments for 'slaveof' command\r\n")
		} else {
			out = []byte("+OK\r\n")
		}
	default:
		out = []byte("-ERR unknown command `" + cmd.Name() + "`, with args beginning with:\r\n")
	}
	return
}

func (s *RedisServer) OnClose(c *connection.Connection) {
	s.log.WithFields(logrus.Fields{
		"action": "Closed",
		"addr":   c.PeerAddr(),
	}).Println()
}
