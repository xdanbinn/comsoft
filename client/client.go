package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}

func main() {
	client := NewClient("127.0.0.1", 8000)
	if client == nil {
		fmt.Println(">>>>>>连接服务器失败....")
		return
	}
	fmt.Println(">>>>>>成功连接服务器....")
	select {}
}
