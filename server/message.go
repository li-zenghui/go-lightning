package server

type Message struct {
	Method    int    `json:"Method"`              // 0 全服  1 群组   2 私聊
	GroupID   string `json:"GroupID,omitempty"`   //群聊房间ID
	OppOpenID string `json:"OppOpenID,omitempty"` //私聊对方ID
	Msg       string `json:"Msg"`                 //消息体
}

// {"Method":0,"Msg":"哈啊哈"}
// {"Method":2,"OppOpenID":"127.0.0.1:57652","Msg":"你也好,007"}
// {"Method":1,"GroupID":"001","Msg":"欢迎大家加入群聊"}
