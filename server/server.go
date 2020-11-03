package server

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"im/proto"
	"im/utils"
	"net"
)

const (
	ADDR = "localhost:8000"
	TYPE = "tcp"
)

func Run() {
	addr, err := net.ResolveTCPAddr(TYPE, ADDR)
	utils.Must(err)
	l,err:=net.ListenTCP("tcp",addr)
	defer l.Close()
	for {
		conn, err := l.AcceptTCP()
		//conn.SetReadBuffer(1)
		//conn.SetReadBuffer()
		//conn.
		utils.Must(err)
		go handleReq(conn)
	}
}

func handleReq(conn *net.TCPConn) {
	//buf := make([]byte, 1024)
	reader:=bufio.NewReader(conn)
	var length int32
	for {
		peek,err:=reader.Peek(4)
		utils.Must(err)
		lenBuf:=bytes.NewBuffer(peek)
		err=binary.Read(lenBuf,binary.BigEndian,&length)
		utils.Must(err)
		if int32(reader.Buffered())<length+4{
			continue
		}
		buf:=make([]byte,length+4)
		_, err = reader.Read(buf)
		//fmt.Println(len)
		//fmt.Println(buf[0:len])
		//fmt.Printf("req len:%d data:%s", len, string(buf))
		utils.Must(err)
		msg:=im.Msg{}
		err=proto.Unmarshal(buf[4:],&msg)
		utils.Must(err)
		utils.PrintStrcut(msg)
		//回复ACK
		msgAck:=im.MsgAck{MsgId:msg.MsgId}
		res,err:=proto.Marshal(&msgAck)
		utils.Must(err)
		writeBuffer,err:=utils.Encode(res)
		utils.Must(err)
		conn.Write(writeBuffer)
	}
	//conn.Close()
}
