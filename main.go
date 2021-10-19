package main

import (
	"fmt"
	"go-lightning/server"
)

func main() {
	S := server.NewServer("127.0.0.1", 8125)
	S.Start()
	fmt.Println("[go-lightning]start,port:8125")
}

// 上线 加载相关group 持久维持group连接
