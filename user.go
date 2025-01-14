package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// 创建一个user的API
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}

	go user.ListenChannel()

	return user
}

// 监听User Channel的方法
func (u *User) ListenChannel() {
	for {
		message := <-u.C
		u.conn.Write([]byte(message + "\n"))
	}
}
