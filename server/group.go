package server

type group struct {
	GroupID string
	Users   map[string]*user
	C       chan string
	Server  *server
}

func NewGroup() *group {
	g := &group{
		GroupID: "001",
		Users:   make(map[string]*user),
		C:       make(chan string),
	}
	go g.ListenGroupMessage()
	return g
}

func (g *group) ListenGroupMessage() {
	for {
		msg := <-g.C
		g.SendAllMessage(msg)
	}
}

func (g *group) SendAllMessage(msg string) {
	for _, u := range g.Users {
		if u != nil {
			u.C <- msg
		}
	}
}

// 加载分组数据
