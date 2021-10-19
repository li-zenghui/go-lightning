package server

import (
	"encoding/json"
	"fmt"
	"net"
)

type user struct {
	OpenID string
	Addr   string
	C      chan string
	Conn   net.Conn
}

func NewUser(conn net.Conn) *user {
	addr := conn.RemoteAddr().String()
	u := &user{
		OpenID: addr,
		Addr:   addr,
		C:      make(chan string),
		Conn:   conn,
	}
	go u.ListenMessage2User()
	return u
}

//监听消息发送到user的消息

func (u *user) ListenMessage2User() {
	for {
		msg := <-u.C
		u.Conn.Write([]byte(msg + "\n"))
	}
}

//处理当前用户发送的消息

func (u *user) HandlerUserMessage(msg []byte, s *server) {
	var message Message
	json.Unmarshal(msg, &message)
	s1 := string(msg)
	fmt.Println(s1)
	switch message.Method {
	case "0":
		//全服信息
		s.BroadCast(u, message.Msg)
	case "1":
		// 群聊
		s.Mu.RLock()
		defer s.Mu.RUnlock()
		g := s.GroupMap[message.GroupID]
		sendData := fmt.Sprintf("[群聊:%v-用户:%v]说:%v", g.GroupID, u.OpenID, message.Msg)
		g.C <- sendData
	case "2":
		//私聊
		s.Mu.RLock()
		defer s.Mu.RUnlock()
		sendData := fmt.Sprintf("[私信-用户:%v]说:%v", u.OpenID, message.Msg)
		s.OnlineMap[message.OppOpenID].C <- sendData
	default:
		s.BroadCast(u, string(msg))
	}
}

//创建群聊
