package server

import (
	"net"
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
	u.server.BroadCast(u, msg)
}

func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg + "\n"))
	}
}
