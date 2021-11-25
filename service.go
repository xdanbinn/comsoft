package main

import (
	"fmt"
)
import "net"

type service struct {
	Ip   string
	Port int
}

func NewService(ip string, port int) *service {
	return &service{
		Ip:   ip,
		Port: port,
	}
}

var i = 1

func (s *service) Handler(con net.Conn) {
	fmt.Println("这是第次test", i)
	i++
}

func (s *service) StartUp() {
	lister, errd := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if errd != nil {
		fmt.Println("net.Listen error:", errd)
	}
	defer lister.Close()
	for {
		conne, err := lister.Accept()
		if err != nil {
			fmt.Println("lister.Accept error:", err)
		}
		go s.Handler(conne)
	}
}
