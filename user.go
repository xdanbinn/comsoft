package main

import "net"

type user struct {
	Name string
	Addr string
	C    chan string
	Con  net.Conn
}

func (s *user) ListenMsg() {
	for {
		msg := <-s.C
		s.Con.Write([]byte(msg + "\n"))
	}
}

func NewUser(conn net.Conn) *user {
	s := &user{
		Name: conn.RemoteAddr().String(),
		Addr: conn.RemoteAddr().String(),
		C:    make(chan string),
		Con:  conn,
	}
	go s.ListenMsg()
	return s
}
