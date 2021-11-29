package main

import (
	"net"
)

type user struct {
	Name string
	Addr string
	C    chan string
	Con  net.Conn

	Ser *Server
}

func (s *user) ListenMsg() {
	for {
		msg := <-s.C
		s.Con.Write([]byte(msg + "\n"))
	}
}

func (s *user) Online() {
	s.Ser.Mux.Lock()
	s.Ser.ClientMap[s.Name] = s
	s.Ser.Mux.Unlock()
	s.Ser.Boadcast(s, "上线了")
}

func (s *user) Offline() {
	s.Ser.Mux.Lock()
	delete(s.Ser.ClientMap, s.Name)
	s.Ser.Mux.Unlock()
	s.Ser.Boadcast(s, "下线了")
}
func (s *user) DoMessage(msg string) {
	if msg == "who" {
		for _, user := range s.Ser.ClientMap {
			msg := "[" + user.Name + "]" + user.Addr + "在线...."
			s.Con.Write([]byte(msg))
		}
	} else {
		s.Ser.Boadcast(s, msg)
	}

}

func NewUser(conn net.Conn, ser *Server) *user {
	s := &user{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		C:    make(chan string),
		Con:  conn,
		Ser:  ser,
	}
	go s.ListenMsg()
	return s
}
