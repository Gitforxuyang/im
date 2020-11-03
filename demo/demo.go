package main

import "fmt"

func main(){
	//msg := im.Msg{}
	//msg.MsgId = strconv.Itoa(int(time.Now().Unix()))
	//msg.Type = im.MsgType_LOGIN
	//msg.Seq = 1
	//var err error
	//msg.Content,err=utils.StructToJson(im.LoginContent{Name:"admin",Password:"admin"})
	//utils.Must(err)
	//buf,err:=proto.Marshal(&msg)
	//utils.Must(err)
	//
	//proto.Unmarshal(buf)
	buf:=[]byte{1,1,11,0,0}
	fmt.Println(buf[0:3])
}



