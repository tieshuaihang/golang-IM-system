package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户的列表
	OnLineMap map[string]*User
	MapLock   sync.RWMutex

	//消息广播的channel
	Message chan string
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnLineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.MapLock.RLock()
		for _, u := range s.OnLineMap {
			u.C <- msg
		}
		s.MapLock.RUnlock()
	}
}

func (s *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + ":" + msg
	s.Message <- sendMsg
}

func (s *Server) Handler(conn net.Conn) {

	u := NewUser(conn, s)

	// 用户上线
	u.OnLine()

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 512)
		for {
			n, err := u.Conn.Read(buf)
			if n == 0 {
				u.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn read err:", err)
				return
			}

			//提取用户的消息(去除'\n')
			msg := string(buf[:n-1])

			u.DoMessage(msg)

		}

	}()

	select {}

}

//启动服务器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	//close listen socket
	defer listener.Close()

	//启动监听Message的goroutine
	go s.ListenMessage()

	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		//do handler
		go s.Handler(conn)
	}
}
