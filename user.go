package main

import (
	"fmt"
	"net"
	"strings"
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
func (s *user) sendMsg(msg string) {
	s.Con.Write([]byte(msg))
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
			s.sendMsg(msg)
		}
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		if _, ok := s.Ser.ClientMap[newName]; ok {
			fmt.Println("用户名已经存在")
		} else {
			s.Ser.Mux.Lock()
			delete(s.Ser.ClientMap, s.Name)
			s.Ser.ClientMap[newName] = s
			s.Ser.Mux.Unlock()
			s.Name = newName
			msg := "您已更新用户名为：" + newName + "\n"
			s.sendMsg(msg)
		}

	} else if strings.HasPrefix(msg, "to|") {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			fmt.Println("请输入对方用户名，to|张三|你好啊！")
			return
		}
		remoteUser, ok := s.Ser.ClientMap[remoteName]
		if !ok {
			fmt.Println("请核对用户名！")
			return
		}
		contx := strings.Split(msg, "|")[2]
		if contx == "" {
			fmt.Println("不可发送空白信息")
			return
		}
		remoteUser.sendMsg(s.Name + "对你说：" + contx)
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
