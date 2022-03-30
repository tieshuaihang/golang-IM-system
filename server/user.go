package server

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	C    chan string
	Conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, s *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		Conn:   conn,
		server: s,
	}
	// 启动监听该用户接受消息goroutine
	go user.ListenMessage()

	return user
}

func (u *User) OnLine() {

	u.server.MapLock.Lock()
	u.server.OnLineMap[u.Name] = u
	u.server.MapLock.Unlock()

	u.server.BroadCast(u, "上线")

}

func (u *User) OffLine() {

	u.server.MapLock.Lock()
	delete(u.server.OnLineMap, u.Name)
	u.server.MapLock.Unlock()

	u.server.BroadCast(u, "下线")

}

func (u *User) DoMessage(msg string) {

	if msg == "who" {
		var onlineMsg string
		u.server.MapLock.RLock()
		for _, u := range u.server.OnLineMap {
			onlineMsg += u.Name + "---"
		}
		onlineMsg += "在线\n"
		u.server.MapLock.RUnlock()
		u.SendMsg(onlineMsg)

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := u.server.OnLineMap[newName]
		if !ok {

			u.server.MapLock.Lock()
			delete(u.server.OnLineMap, u.Name)
			u.server.OnLineMap[newName] = u
			u.server.MapLock.Unlock()
			u.Name = newName
			u.SendMsg("修改用户名成功\n")
		} else {
			u.SendMsg("该用户名已经被使用\n")
		}

	} else {
		u.server.BroadCast(u, msg)
	}

}

//给当前User对应的客户端发送消息
func (u *User) SendMsg(msg string) {
	u.Conn.Write([]byte(msg))
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg + "\n"))
	}
}
