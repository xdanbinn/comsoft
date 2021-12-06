package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	Message   chan string
	ClientMap map[string]*user
	Mux       sync.RWMutex
}

func NewService(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		Message:   make(chan string),
		ClientMap: make(map[string]*user),
	}
}
func (s *Server) ListMsg() {
	for {
		msg := <-s.Message
		s.Mux.Lock()
		for _, key := range s.ClientMap {
			key.C <- msg
		}
		s.Mux.Unlock()
	}

}

func (s *Server) Boadcast(user *user, msg string) {
	sendMsg := "[" + user.Name + "]" + user.Addr + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(con net.Conn) {
	user := NewUser(con, s)
	user.Online()
	timeout := make(chan bool)
	go func() {
		mes := make([]byte, 4096)
		for {
			n, err := con.Read(mes)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read Message error:", err)
				return
			}
			msg := string(mes[0 : n-1])
			user.DoMessage(msg)
			timeout <- true
		}
	}()
	for {
		select {
		case <-timeout:
		case <-time.After(time.Second * 600):
			user.sendMsg("与服务器连接超时....")
			close(user.C)
			con.Close()
			return
		}
	}

}

func (s *Server) StartUp() {
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
