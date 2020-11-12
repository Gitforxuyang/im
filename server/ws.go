package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/gorilla/websocket"
	"im/proto/protocol"
	"im/utils"
	"net/http"
)

const (
	addr = "localhost:8080"
)

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	utils.Must(err)
	fmt.Printf("connection accept remote:%s", c.RemoteAddr().String())
	utils.Must(err)
	defer c.Close()
	for {
		msgType, data, err := c.ReadMessage()
		utils.Must(err)
		fmt.Println(msgType)
		fmt.Println(len(data))
		fmt.Println(data)
		fmt.Printf("data:%s \n", string(data))
		msg := protocol.NimProtocol{}
		var writeBuf = new(bytes.Buffer)
		err = binary.Write(writeBuf, binary.BigEndian, msg.CmdId)
		utils.Must(err)
		err = binary.Write(writeBuf, binary.BigEndian, msg.Version)
		utils.Must(err)
		err = binary.Write(writeBuf, binary.BigEndian, msg.BodyLen)
		utils.Must(err)
		err = binary.Write(writeBuf, binary.BigEndian, msg.Body)
		utils.Must(err)
		err = c.WriteMessage(websocket.BinaryMessage, writeBuf.Bytes())
		utils.Must(err)
	}
}

func Ws() {
	http.HandleFunc("/echo", echo)
	http.ListenAndServe(addr, nil)

}
