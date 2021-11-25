package main

import (
	"fmt"
	"net"
	"sync"
)

type service struct {
	Ip        string
	Port      int
	Message   chan string
	ClientMap map[string]*user
	Mux       sync.RWMutex
}

func NewService(ip string, port int) *service {
	return &service{
		Ip:        ip,
		Port:      port,
		Message:   make(chan string),
		ClientMap: make(map[string]*user),
	}
}
func (s *service) ListMsg() {
	for {
		msg := <-s.Message
		s.Mux.Lock()
		for _, key := range s.ClientMap {
			key.C <- msg
		}
		s.Mux.Unlock()
	}

}

func (s *service) Boadcast(user *user, msg string) {
	sendMsg := "[" + user.Name + "]" + user.Addr + msg
	s.Message <- sendMsg

}

func (s *service) Handler(con net.Conn) {
	//新用户建立连接了
	user := NewUser(con)
	//服务器注册了
	s.Mux.Lock()
	s.ClientMap[user.Name] = user
	s.Mux.Unlock()
	//服务器广播
	s.Boadcast(user, "上线了")

	select {}

}

func (s *service) StartUp() {
	lister, errd := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if errd != nil {
		fmt.Println("net.Listen error:", errd)
	}
	defer lister.Close()
	go s.ListMsg()
	for {
		conne, err := lister.Accept()
		if err != nil {
			fmt.Println("lister.Accept error:", err)
		}
		go s.Handler(conne)
	}
}
