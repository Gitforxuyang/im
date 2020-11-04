package conn

import "im/proto"

type Msg struct {
	msg *im.Msg
	//conn Conn
}

func (m *Msg) Ack() error{
	return nil
}
