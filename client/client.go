package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	"im/proto"
	"im/utils"
	"math/rand"
	"net"
	"strconv"
	"time"
)

func main() {
	client()
}

func client() {
	addr := "localhost:8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	utils.Must(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	utils.Must(err)
	var seq int32 = 1
	for {
		time.Sleep(time.Second * time.Duration(rand.Intn(20)))
		msg := im.Msg{}
		msg.MsgId = strconv.Itoa(int(time.Now().Unix()))
		msg.Type = im.MsgType_LOGIN
		seq++
		msg.Seq = seq
		msg.Content,err=proto.Marshal(&im.LoginContent{Name:`admin`,Password:"admin"})
		buf,err:=proto.Marshal(&msg)
		utils.Must(err)
		buf,err=utils.Encode(buf)
		utils.Must(err)
		_,err=conn.Write(buf)
		utils.Must(err)
		reader:=bufio.NewReader(conn)

		var length int32
		peek,err:=reader.Peek(4)
		utils.Must(err)
		lenBuf:=bytes.NewBuffer(peek)
		err=binary.Read(lenBuf,binary.BigEndian,&length)
		utils.Must(err)
		res:=make([]byte,length+4)
		_, err = reader.Read(res)

		msgAck:=im.MsgAck{}
		err=proto.Unmarshal(res[4:],&msgAck)
		utils.Must(err)
		utils.PrintStrcut(msgAck)
		time.Sleep(time.Second*1000)
		//return
	}
}
