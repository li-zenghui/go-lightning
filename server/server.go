package server

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*user
	GroupMap  map[string]*group
	Mu        sync.RWMutex
	//消息广播信道
	C chan string
}

func NewServer(Ip string, port int) *server {
	return &server{
		Ip:        Ip,
		Port:      port,
		OnlineMap: make(map[string]*user),
		GroupMap:  make(map[string]*group),
		C:         make(chan string),
	}
}

func (s *server) GetUrl() string {
	return fmt.Sprintf(":%v", s.Port)
}

func (s *server) Start() {
	//监听端口
	listen, err := net.Listen("tcp", s.GetUrl())
	if err != nil {
		fmt.Printf("net.listen err,%v", err)
	}
	defer listen.Close()
	go s.ListenMsg()
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("net.accept err,%v", err)
			continue
		}
		go s.Handler(conn)
	}
}

func (s *server) Handler(conn net.Conn) {
	user := NewUser(conn)
	s.Mu.Lock()
	s.OnlineMap[user.OpenID] = user
	s.Mu.Unlock()
	s.BroadCast(user, "我上线了")
	// 监听用户消息
	go s.ListenUserMessage(user)
	select {}
}

func (s *server) ListenUserMessage(user *user) {
	msg := make([]byte, 4096)
	for {
		n, err := user.Conn.Read(msg)
		if err != nil && err != io.EOF {
			fmt.Printf("io.read err,%v", err)
			return
		}
		if n == 0 {
			s.BroadCast(user, "我下线了")
			return
		}
		user.HandlerUserMessage(msg[0:n], s)
	}

}

func (s *server) BroadCast(user *user, msg string) {
	sendMsg := fmt.Sprintf("[全服-用户:%v]说:%v", user.OpenID, msg)
	s.C <- sendMsg
}

func (s *server) ListenMsg() {
	for {
		msg := <-s.C
		s.Mu.RLock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.Mu.RUnlock()
	}
}
