package main

import "go-lightning/server"

func main() {
	S := server.NewServer("127.0.0.1", 8123)
	S.Start()
}

// 上线 加载相关group 持久维持group连接
