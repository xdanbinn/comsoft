package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}
	client.conn = conn
	return client
}
func (c *Client) menu() bool {

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.重命名模式")
	fmt.Println("0.退出")
	var falg int
	fmt.Scan(&falg)
	if falg >= 0 && falg <= 3 {
		c.flag = falg
		return true
	} else {
		fmt.Println(">>>>请输入合适的数字<<<<")
		return false
	}
}
func (c *Client) Run() {
	for c.flag != 0 {
		for c.menu() != true {

		}
		switch c.flag {
		case 1:
			fmt.Println("选择公聊模式...")
		case 2:
			fmt.Println("选择私聊模式...")
		case 3:
			fmt.Println("选择重命名模式...")
		case 0:
			fmt.Println("退出...")
		}
	}

}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "I", "127.0.0.1", "设置服务器地址")
	flag.IntVar(&serverPort, "P", 8000, "设置服务器端口")
}

func main() {
	flag.Parse()
	//client := NewClient("127.0.0.1", 8000)
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>连接服务器失败....")
		return
	}
	fmt.Println(">>>>>>成功连接服务器....")
	client.Run()
}
