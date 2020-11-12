package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"im/conn/tcp"
	"im/proto"
	"im/proto/protocol"
	"im/utils"
	"net"
)

func main() {
	client()
}

func client() {
	addr := "localhost:8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	utils.Must(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	c:=tcp.NewTCPConn(conn)
	msg := im.AuthReq{}
	msg.Uid = "123"
	msg.Token = "token"
	p := protocol.NimProtocol{}
	p.Version = 1
	p.CmdId = protocol.AuthReq
	p.Body, _ = proto.Marshal(&msg)
	p.BodyLen = uint16(len(p.Body))
	c.Write(&p)
	res,err:=c.Read()
	utils.Must(err)
	resProto:=im.AuthRes{}
	err=proto.Unmarshal(res.Body,&resProto)
	utils.Must(err)
	fmt.Println(resProto)
}

