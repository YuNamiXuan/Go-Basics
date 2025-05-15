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

	// close(u.C)
	// 广播用户下线消息
	u.server.BroadCast(u, "已下线")
}

// 给当前用户的客户端发送消息
func (u *User) SendMessage(message string) {
	u.conn.Write([]byte(message))
}

// 处理消息
func (u *User) ProcessMessage(message string) {
	// message = strings.TrimSpace(message)

	if message == "who" {
		// 使用who查询在线用户列表
		u.server.maplock.Lock()
		for _, user := range u.server.OnlineMap {
			onlineMessage := "[" + user.Addr + "]" + user.Name + "在线...\n"
			u.SendMessage(onlineMessage)
		}
		u.server.maplock.Unlock()

	} else if strings.HasPrefix(message, "rename|") {
		// 使用rename|修改用户名
		parts := strings.SplitN(message, "|", 2)
		if len(parts) != 2 {
			u.SendMessage("格式错误，正确格式为: rename|新用户名\n")
			return
		}
		newName := parts[1]
		u.server.maplock.Lock()
		if _, exists := u.server.OnlineMap[newName]; exists {
			u.SendMessage("该用户已存在\n")
		} else {
			delete(u.server.OnlineMap, u.Name)
			u.Name = newName
			u.server.OnlineMap[newName] = u
			u.SendMessage("用户名修改成功\n")
		}
		u.server.maplock.Unlock()

	} else if len(message) > 4 && message[:3] == "to|" {
		// 单独发送用户消息 格式：to|username|message
		parts := strings.SplitN(message, "|", 3)
		if len(parts) != 3 {
			u.SendMessage("格式错误，正确格式为: to|用户名|消息内容\n")
		}
		username, sendMessage := parts[1], parts[2]

		u.server.maplock.Lock()
		user, exist := u.server.OnlineMap[username]
		u.server.maplock.Unlock()
		if !exist {
			u.SendMessage("该用户不存在\n")
			return
		}
		user.SendMessage(u.Name + "对你说: " + sendMessage)

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
