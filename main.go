package main

import "im/server"

func main() {
	//handler := handler2.NewHandler()
	//
	//server := tcp.NewTcpServer(8000, handler)
	//server.Run()
	server.Ws()
}
