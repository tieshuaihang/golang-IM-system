package server

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	isAlive := make(chan struct{})

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
			msg := string(buf[:n])

			u.DoMessage(msg)

			// 用户的任何消息，都说明其在线
			isAlive <- struct{}{}

		}

	}()

	for {
		select {
		case <-isAlive:
			//当前用户是活跃的，应该重置定时器
			//不做任何事情，为了激活select，更新下面的定时器
		case <-time.After(100 * time.Second):
			// 用户超时
			u.SendMsg("您超时了\n")
			//关闭资源
			close(u.C)

			conn.Close()
			return
		}

	}

}

//启动服务器的接口
func (s *Server) Start() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			return
		}
	}()

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
