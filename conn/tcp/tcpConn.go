package tcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"im/conn"
	"im/proto/protocol"
	"io"
	"net"
)

const (
	READ_BUFFER_SIZE  = 8192  //读取缓冲区大小
	WRITE_BUFFER_SIZE = 1024  //写入缓冲区大小
	MSG_MAX_SIZE      = 10240 //最大消息大小，防止被人攻击
)

type tcpConn struct {
	conn       *net.TCPConn
	reader     *bufio.Reader
	writer     *bufio.Writer
	remoteAddr string
	//最新更新时间
	pingedAt int64
	login    bool //是否已登陆
	uid      int64
}

//5分钟过期
func (m *tcpConn) GetExpireAt() int64 {
	return m.pingedAt + 300000
}

func (m *tcpConn) SetUid(uid int64) {
	m.uid = uid
}

func (m *tcpConn) GetUid() int64 {
	return m.uid
}

func (m *tcpConn) SetLogin(login bool) {
	m.login = login
}

func (m *tcpConn) GetLogin() bool {
	return m.login
}

func (m *tcpConn) Close() {
	m.conn.Close()
}

func (m *tcpConn) GetPingedAt() int64 {
	return m.pingedAt
}

func (m *tcpConn) SetPingedAt(time int64) {
	m.pingedAt = time
}

func (m *tcpConn) RemoteAddr() string {
	return m.remoteAddr
}

func (m *tcpConn) Read() (*protocol.NimProtocol, error) {
	headerBuf := make([]byte, 4)
	err := binary.Read(m.reader, binary.BigEndian, headerBuf)
	if err != nil {
		if err == io.EOF {
			return nil, errors.New("connection closed by client")
		}
		if err.Error() != "use of closed network connection" {
			return nil, errors.New("connection closed by server")
		}
		return nil, err
	}
	p := &protocol.NimProtocol{}
	p.CmdId = protocol.CmdEnum(uint8(headerBuf[0]))
	p.Version = uint8(headerBuf[1])
	p.BodyLen = binary.BigEndian.Uint16(headerBuf[2:4])
	if p.BodyLen != 0 {
		bodyBuf := make([]byte, p.BodyLen)
		err = binary.Read(m.reader, binary.BigEndian, bodyBuf)
		if err != nil {
			if err.Error() != "use of closed network connection" {
				return nil, errors.New("connection closed by server")
			}
			return nil, err
		}
		p.Body = bodyBuf
	}
	return p, nil
}

func (m *tcpConn) Write(msg *protocol.NimProtocol) error {
	var writeBuf = new(bytes.Buffer)
	err := binary.Write(writeBuf, binary.BigEndian, msg.CmdId)
	if err != nil {
		return err
	}
	err = binary.Write(writeBuf, binary.BigEndian, msg.Version)
	if err != nil {
		return err
	}
	err = binary.Write(writeBuf, binary.BigEndian, msg.BodyLen)
	if err != nil {
		return err
	}
	err = binary.Write(writeBuf, binary.BigEndian, msg.Body)
	if err != nil {
		return err
	}
	_, err = m.writer.Write(writeBuf.Bytes())
	if err != nil {
		if err.Error() != "use of closed network connection" {
			return errors.New("connection closed by server")
		}
		return err
	}
	m.writer.Flush()
	if err != nil {
		if err.Error() != "use of closed network connection" {
			return errors.New("connection closed by server")
		}
		return err
	}
	return nil
}

func NewTCPConn(conn *net.TCPConn) conn.Conn {
	conn.SetReadBuffer(READ_BUFFER_SIZE)
	conn.SetWriteBuffer(WRITE_BUFFER_SIZE)
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	return &tcpConn{reader: reader, writer: writer, remoteAddr: conn.RemoteAddr().String(), conn: conn}
}
