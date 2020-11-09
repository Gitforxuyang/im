package protocol

import (
	"github.com/golang/protobuf/proto"
	"im/proto"
)

//自定义协议
type NimProtocol struct {
	CmdId   CmdEnum //数据类型
	Version uint8   //版本号 目前默认 1
	BodyLen uint16  //数据包大小  最大65535  60K多
	Body    []byte  //主体
}

type CmdEnum uint8

const (
	Ping       CmdEnum = 1
	Pong       CmdEnum = 2
	AuthReq    CmdEnum = 3 //认证请求
	AuthRes    CmdEnum = 4 //认证结果
	LogoutReq  CmdEnum = 5 //退出请求
	LogoutRes  CmdEnum = 6 //退出结果
	KickoutReq CmdEnum = 7 //踢人
	//KickoutRes     CmdEnum = 8
	PullMsgReq     CmdEnum = 9 //拉取消息
	PullMsgRes     CmdEnum = 10
	MessageSend    CmdEnum = 11 //消息
	MessageSendAck CmdEnum = 12 //发送消息ACK
	MessageNotify  CmdEnum = 13 //下发聊天消息
	Notify         CmdEnum = 14 //通知类消息 下发通知、派对被关闭、帖子被删等消息
	NotifyAck      CmdEnum = 15 // 通知类的ACK
)

func MakePong() *NimProtocol {
	protocol := NimProtocol{}
	protocol.CmdId = Pong
	protocol.BodyLen = 0
	protocol.Version = 1
	return &protocol
}

func MakeAuthRes(msg *im.AuthRes) *NimProtocol {
	protocol := NimProtocol{}
	protocol.CmdId = AuthRes
	protocol.Version = 1
	protocol.Body, _ = proto.Marshal(msg)
	protocol.BodyLen = uint16(len(protocol.Body))
	return &protocol
}

func MakeLogoutRes(msg *im.LogoutRes) *NimProtocol {
	protocol := NimProtocol{}
	protocol.CmdId = LogoutRes
	protocol.Version = 1
	protocol.Body, _ = proto.Marshal(msg)
	protocol.BodyLen = uint16(len(protocol.Body))
	return &protocol
}

func MakeKickoutReq(msg *im.KickoutReq) *NimProtocol {
	protocol := NimProtocol{}
	protocol.CmdId = KickoutReq
	protocol.Version = 1
	protocol.Body, _ = proto.Marshal(msg)
	protocol.BodyLen = uint16(len(protocol.Body))
	return &protocol
}

func MakeMessageSendAck(msg *im.MessageSendAck) *NimProtocol {
	protocol := NimProtocol{}
	protocol.CmdId = MessageSendAck
	protocol.Version = 1
	protocol.Body, _ = proto.Marshal(msg)
	protocol.BodyLen = uint16(len(protocol.Body))
	return &protocol
}


func MakeMessageNotify(msg *im.MessageNotify) *NimProtocol {
	protocol := NimProtocol{}
	protocol.CmdId = MessageNotify
	protocol.Version = 1
	protocol.Body, _ = proto.Marshal(msg)
	protocol.BodyLen = uint16(len(protocol.Body))
	return &protocol
}


func MakePullMessageRes(msg *im.PullMsgRes) *NimProtocol {
	protocol := NimProtocol{}
	protocol.CmdId = PullMsgRes
	protocol.Version = 1
	protocol.Body, _ = proto.Marshal(msg)
	protocol.BodyLen = uint16(len(protocol.Body))
	return &protocol
}

func Unmarshal(msg *NimProtocol) (proto.Message, error) {
	var err error
	switch msg.CmdId {
	case AuthReq:
		message := im.AuthReq{}
		err = proto.Unmarshal(msg.Body, &message)
		if err != nil {
			return nil, err
		}
		return &message, nil
	case LogoutReq:
		message := im.LogoutReq{}
		err = proto.Unmarshal(msg.Body, &message)
		if err != nil {
			return nil, err
		}
		return &message, nil
	case PullMsgReq:
		message := im.PullMsgReq{}
		err = proto.Unmarshal(msg.Body, &message)
		if err != nil {
			return nil, err
		}
		return &message, nil
	case MessageSend:
		message := im.MessageSend{}
		err = proto.Unmarshal(msg.Body, &message)
		if err != nil {
			return nil, err
		}
		return &message, nil
	case NotifyAck:
		message := im.NotifyAck{}
		err = proto.Unmarshal(msg.Body, &message)
		if err != nil {
			return nil, err
		}
		return &message, nil
	}
	return nil, nil
}
