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
	Ping CmdEnum = 1
	Pong CmdEnum = 2
	Ack  CmdEnum = 3 //消息ACK
	Message  CmdEnum = 4 //消息
	//ChatRoom   CmdEnum = 4  //聊天室
	//GroupChat  CmdEnum = 5  //群聊
	//SingleChat CmdEnum = 6  //单聊
	AuthReq    CmdEnum = 7  //认证请求
	AuthRes    CmdEnum = 8  //认证结果
	LogoutReq  CmdEnum = 9  //退出请求
	LogoutRes  CmdEnum = 10 //退出结果
	KickoutReq CmdEnum = 11 //踢人
	KickoutRes CmdEnum = 12
	PullMsgReq CmdEnum = 13 //拉取消息
	PullMsgRes CmdEnum = 14
	PushMsgReq CmdEnum = 15 //推送消息
	PushMsgRes CmdEnum = 16
	Notify     CmdEnum = 17 //通知类消息 下发通知、派对被关闭、帖子被删等消息
)
