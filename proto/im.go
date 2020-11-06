package im

//自定义协议
type NimProtocol struct {
	CmdId   uint8  //数据类型
	Version uint8  //版本号 目前默认 1
	BodyLen uint16 //数据包大小  最大65535  60K多
	Body    []byte //主体
}

type CmdEnum uint8

const (
	Ping           CmdEnum = 1
	Pong           CmdEnum = 2
	AuthReq        CmdEnum = 3 //认证请求
	AuthRes        CmdEnum = 4 //认证结果
	LogoutReq      CmdEnum = 5 //退出请求
	LogoutRes      CmdEnum = 6 //退出结果
	KickoutReq     CmdEnum = 7 //踢人
	KickoutRes     CmdEnum = 8
	PullMsgReq     CmdEnum = 9 //拉取消息
	PullMsgRes     CmdEnum = 10
	MessageSend    CmdEnum = 11 //消息
	MessageSendAck CmdEnum = 12 //发送消息ACK
	MessageNotify  CmdEnum = 13 //下发聊天消息
	Notify         CmdEnum = 14 //通知类消息 下发通知、派对被关闭、帖子被删等消息
	NotifyAck      CmdEnum = 15 // 通知类的ACK
)
