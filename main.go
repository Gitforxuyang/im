package main

import (
	"im/conn/tcp"
	handler2 "im/handler"
)

func main() {
	//server.Run()
	handler := handler2.NewHandler()

	server := tcp.NewTcpServer(8000, handler)
	server.Run()
}
