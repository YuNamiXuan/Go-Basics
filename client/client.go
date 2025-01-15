package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

// 创建client
func NewClient(ip string, port int) *Client {
	c := &Client{
		ServerIp:   ip,
		ServerPort: port,
		flag:       9999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Println("Error net.Dial:", err)
		return nil
	}
	c.conn = conn
	return c
}

func (c *Client) ProcessResponse() {
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1.公开聊天")
	fmt.Println("2.单独聊天")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围的数字")
		return false
	}
}

// 查询在线用户
func (c *Client) SelectUsers() {
	sendMessage := "who\n"
	_, err := c.conn.Write([]byte(sendMessage))
	if err != nil {
		fmt.Println("Error conn.Write:", err)
		return
	}
}

func (c *Client) PrivateChat() {
	var chatName string
	var chatMessage string

	c.SelectUsers()
	fmt.Println("请输入要发送的用户名(输入exit退出):")
	fmt.Scanln(&chatName)

	for chatName != "exit" {
		fmt.Println("请输入要发送的消息:(输入exit退出):")
		fmt.Scanln(&chatMessage)

		for chatMessage != "exit" {
			if len(chatMessage) != 0 {
				sendMessage := "to|" + chatName + "|" + chatMessage + "\n"
				_, err := c.conn.Write([]byte(sendMessage))
				if err != nil {
					fmt.Println("Error conn.Write:", err)
					break
				}
			}

			chatMessage = ""
			fmt.Println("请输入要发送的消息:(输入exit退出):")
			fmt.Scanln(&chatMessage)
		}

		c.SelectUsers()
		fmt.Println("请输入要发送的用户名(输入exit退出):")
		fmt.Scanln(&chatName)
	}
}

func (c *Client) ReName() bool {
	fmt.Println("输入更改后的用户名")
	fmt.Scanln(&c.Name)

	sendMessage := "rename|" + c.Name + "\n"
	_, err := c.conn.Write([]byte(sendMessage))
	if err != nil {
		fmt.Println("Error conn.Write:", err)
		return false
	}
	return true
}

func (c *Client) PublicChat() {
	var chatMessage string
	fmt.Println("请输入要发送的消息:(输入exit退出)")
	fmt.Scanln(&chatMessage)

	for chatMessage != "exit" {
		if len(chatMessage) != 0 {
			sendMessage := chatMessage + "\n"
			_, err := c.conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("Error conn.Write:", err)
				break
			}
		}

		chatMessage = ""
		fmt.Println("请输入要发送的消息:(输入exit退出)")
		fmt.Scanln(&chatMessage)
	}
}

func (c *Client) Run() {
	for c.flag != 0 {
		for !c.menu() {
		}

		switch c.flag {
		case 1:
			c.PublicChat()
		case 2:
			c.PrivateChat()
		case 3:
			c.ReName()
		}
	}
}

var serverIp string
var serverPort int

// 命令行解析
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认为127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口号(默认为8888)")
}
