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

// 创建client
func NewClient(ip string, port int) *Client {
	c := &Client{
		ServerIp:   ip,
		ServerPort: port,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("Error net.Dial:", err)
		return nil
	}
	c.conn = conn
	return c
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("连接服务器失败")
	} else {
		fmt.Println("连接服务器成功")
	}

	select {}
}
