package conn

//连接
type Conn interface {
	//读取消息
	Read() (Msg,error)
	Write(msg Msg) error
}


