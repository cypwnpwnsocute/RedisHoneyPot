// @Title  main.go
// @Description A highly interactive honeypot supporting redis protocol
// @Author  Cy 2021.04.08
package main

func init() {

}

func main() {
	s, err := NewRedisServer("0.0.0.0:6379", "tcp", 1)
	if err != nil {
		panic(err)
	}
	defer s.Stop()
	s.Start()
}
