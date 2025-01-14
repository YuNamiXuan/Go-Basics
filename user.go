package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// 创建user
func NewUser(conn net.Conn, s *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: s,
	}

	go user.ListenChannel()

	return user
}

// 用户上线
func (u *User) Online() {
	// 将用户加入在线列表中
	u.server.maplock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.maplock.Unlock()

	// 广播用户上线消息
	u.server.BroadCast(u, "已上线")
}

// 用户下线
func (u *User) Offline() {
	// 将用户从在线列表中删除
	u.server.maplock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.maplock.Unlock()

	// 广播用户下线消息
	u.server.BroadCast(u, "已下线")
}

// 给当前用户发送消息
func (u *User) SendMessage(message string) {
	u.conn.Write([]byte(message))
}

// 处理消息
func (u *User) ProcessMessage(message string) {

	if message == "who" {
		// 使用who查询在线用户列表
		u.server.maplock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMessage := "[" + user.Addr + "]" + user.Name + "\n"
			u.SendMessage(onlineMessage)
		}
		u.server.maplock.Unlock()

	} else if len(message) > 7 && message[:7] == "rename|" {
		// 使用rename|修改用户名
		newName := strings.Split(message, "|")[1]
		_, ok := u.server.OnlineMap[newName]
		if ok {
			u.SendMessage("该用户名已存在")
		} else {
			u.server.maplock.Lock()
			delete(u.server.OnlineMap, newName)
			u.server.OnlineMap[newName] = u
			u.server.maplock.Unlock()
			u.Name = newName
			u.SendMessage("用户名修改成功")
		}

	} else if len(message) > 4 && message[:3] == "to|" {
		// 单独发送用户消息 格式：to|username|message
		username := strings.Split(message, "|")[1]
		sendmessage := strings.Split(message, "|")[2]
		if username == "" {
			u.SendMessage("消息格式不正确")
			return
		}
		user, ok := u.server.OnlineMap[username]
		if !ok {
			u.SendMessage("该用户不存在")
			return
		}
		if sendmessage == "" {
			u.SendMessage("当前消息为空")
		} else {
			user.SendMessage(u.Name + "说: " + sendmessage)
		}

	} else {
		// 广播用户消息
		u.server.BroadCast(u, message)
	}
}

// 监听User Channel
func (u *User) ListenChannel() {
	for {
		message := <-u.C
		u.conn.Write([]byte(message + "\n"))
	}
}
