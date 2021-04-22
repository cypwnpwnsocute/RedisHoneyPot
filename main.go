// @Title  main.go
// @Description High Interaction Honeypot Solution for Redis protocol
// @Author  Cy 2021.04.08
package main

import "flag"

var (
	addr  string
	proto string
	num   int
)

func init() {
	flag.StringVar(&addr, "addr", "0.0.0.0:6379", "listen address")
	flag.StringVar(&proto, "proto", "tcp", "listen proto")
	flag.IntVar(&num, "num", 1, "loops num")
	flag.Parse()
}

func main() {
	s, err := NewRedisServer(addr, proto, num)
	if err != nil {
		panic(err)
	}
	defer s.Stop()
	s.Start()
}
