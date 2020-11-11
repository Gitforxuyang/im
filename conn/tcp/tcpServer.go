package tcp

import (
	"fmt"
	"im/conn"
	"im/handler"
	"im/utils"
	"net"
)

type tcpServer struct {
	port    int32
	handler *handler.Handler
}

func (m *tcpServer) Run() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", m.port))
	utils.Must(err)
	l, err := net.ListenTCP("tcp", addr)
	utils.Must(err)
	defer l.Close()
	fmt.Printf("tcp server run port:%d  \n", m.port)
	for {
		tcpConn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("server accept error", err)
			tcpConn.Close()
		} else {
			fmt.Printf("client connected to server info:%s  \n", tcpConn.RemoteAddr())
			conn := NewTCPConn(tcpConn)
			go m.handler.Open(conn)
		}
	}
}

func NewTcpServer(port int32, handler *handler.Handler) conn.Server {
	return &tcpServer{port: port, handler: handler}
}
