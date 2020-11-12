package server

import (
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
	l, err := net.ListenTCP("tcp", addr)
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
	//conn.SetReadBuffer(8192)
	//reader := bufio.NewReader(conn)
	////for {
	////	buf := make([]byte, 1024)
	////	fmt.Println(reader.Buffered())
	////	reader.Read(buf)
	////
	////}
	//var length int32
	//for {
	//	peek, err := reader.Peek(4)
	//	utils.Must(err)
	//	fmt.Println(string(peek))
	//	lenBuf := bytes.NewBuffer(peek)
	//	err = binary.Read(lenBuf, binary.BigEndian, &length)
	//	utils.Must(err)
	//	buf := make([]byte, length+4)
	//	err=binary.Read(reader,binary.BigEndian,buf)
	//
	//	utils.Must(err)
	//	msg := im.Msg{}
	//	err = proto.Unmarshal(buf[4:], &msg)
	//	utils.Must(err)
	//	utils.PrintStrcut(msg)
	//	//回复ACK
	//	msgAck := im.MsgAck{MsgId: msg.MsgId}
	//	res, err := proto.Marshal(&msgAck)
	//	utils.Must(err)
	//	writeBuffer, err := utils.Encode(res)
	//	utils.Must(err)
	//	err=binary.Write(conn,binary.BigEndian,writeBuffer)
	//	utils.Must(err)
	//	//conn.Write(writeBuffer)
	//}
	//conn.Close()
}
