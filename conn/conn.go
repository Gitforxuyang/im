package conn

import (
	"im/proto/protocol"
)

//连接
type Conn interface {
	//读取消息
	Read() (*protocol.NimProtocol, error)
	Write(msg *protocol.NimProtocol) error
	RemoteAddr() string
	GetPingedAt() int64
	SetPingedAt(time int64)
	Close()
	SetLogin(login bool)
	GetLogin() bool
	SetUid(uid int64)
	GetUid() int64
	//获取连接的过期时间
	GetExpireAt() int64
}
