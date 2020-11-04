package tcp

import (
	"bufio"
	"im/conn"
	"net"
)
const (
	READ_BUFFER_SIZE=8192 //读取缓冲区大小
	WRITE_BUFFER_SIZE=1024	//写入缓冲区大小
	MSG_MAX_SIZE=10240 //最大消息大小，防止被人攻击
)
type tcpConn struct {
	//conn *net.TCPConn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (m *tcpConn) Read() (conn.Msg,error) {
	panic("implement me")
}

func (m *tcpConn) Write(msg conn.Msg)error {
	panic("implement me")
}

func NewTCPConn(conn *net.TCPConn) conn.Conn{
	conn.SetReadBuffer(READ_BUFFER_SIZE)
	conn.SetWriteBuffer(WRITE_BUFFER_SIZE)
	reader:=bufio.NewReader(conn)
	writer:=bufio.NewWriter(conn)
	return &tcpConn{reader:reader,writer:writer}
}
