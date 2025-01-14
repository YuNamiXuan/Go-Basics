package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	maplock   sync.RWMutex
	Message   chan string
}

// 创建一个server的API
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

// 监听Message Channel，一旦有消息发送至所有user
func (s *Server) ListenMessage() {
	for {
		message := <-s.Message

		s.maplock.Lock()
		for _, u := range s.OnlineMap {
			u.C <- message
		}
		s.maplock.Unlock()
	}
}

// 将上线消息传递至Message Channel
func (s *Server) BroadCast(u *User, m string) {
	sendM := "[" + u.Addr + "]" + u.Name + ":" + m
	s.Message <- sendM
}

func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn)

	// 将用户加入在线列表中
	s.maplock.Lock()
	s.OnlineMap[user.Name] = user
	s.maplock.Unlock()

	// 广播用户上线消息
	s.BroadCast(user, "已上线")

	// 广播用户发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "下线")
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Error read conn:", err)
				return
			}

			message := string(buf[:n-1])
			s.BroadCast(user, message)
		}
	}()

	select {}
}

// 启动服务器的方法
func (s *Server) Start() {
	// socket listening
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}

	defer listener.Close()

	go s.ListenMessage()

	for {
		// accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// do handler
		go s.Handler(conn)

	}
}
