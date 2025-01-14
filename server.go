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
	OnlineMap map[string]*User
	maplock   sync.RWMutex
	Message   chan string
}

// 创建server
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

// 广播消息
func (s *Server) BroadCast(u *User, m string) {
	sendM := "[" + u.Addr + "]" + u.Name + ":" + m
	s.Message <- sendM
}

// 处理user需求
func (s *Server) Handler(conn net.Conn) {
	user := NewUser(conn, s)

	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Error read conn:", err)
				return
			}

			message := string(buf[:n-1])
			user.ProcessMessage(message)

			isLive <- true
		}
	}()

	for {
		select {
		case <-isLive:

		case <-time.After(time.Second * 30):
			user.SendMessage("已超时，被强制下线")
			close(user.C)
			conn.Close()
			return
		}
	}
}

// 启动server
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
