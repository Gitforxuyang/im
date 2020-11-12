package handler

import (
	"fmt"
	"im/conn"
	"im/proto"
	"im/proto/protocol"
	"im/utils"
	"sync"
	"time"
)

type Handler struct {
	//等待验证合法性的连接
	waitAuthSocket sync.Map
	//已经验证完毕的连接
	connectedSocket sync.Map
	//重试队列
	retryList *utils.RetryList
	//msgId对应重试消息
	msgIdMap sync.Map
}

//一个新的连接建立
func (m *Handler) Open(conn conn.Conn) {
	conn.SetPingedAt(utils.NowMillisecond())
	//当连接建立时将连接保存
	m.waitAuthSocket.Store(conn.RemoteAddr(), conn)
	//循环处理
	m.loop(conn)

}
func (m *Handler) loop(conn conn.Conn) {
	for {
		msg, err := conn.Read()
		if err != nil {
			fmt.Println("loop read data error: ", err)
			//当发现读取数据有错误后，关闭连接
			conn.Close()
			m.waitAuthSocket.Delete(conn.RemoteAddr())
			m.connectedSocket.Delete(conn.GetUid())
			return
		}
		var res *protocol.NimProtocol
		switch msg.CmdId {
		case protocol.Ping:
			res, err = m.ping(conn)
		case protocol.AuthReq:
			res, err = m.authReq(conn, msg)
		case protocol.LogoutReq:
			res, err = m.logoutReq(conn, msg)
		case protocol.PullMsgReq:
			res, err = m.pullMsgReq(conn, msg)
		case protocol.MessageSend:
			res, err = m.messageSend(conn, msg)
		case protocol.NotifyAck:
			res, err = m.notifyAck(conn, msg)
		default:
			fmt.Printf("receive unknow cmdId %d", msg.CmdId)
		}
		//如果在处理业务的过程中，遇到了如回写报错等错误，则断掉连接
		if err != nil {
			fmt.Println("loop read data error: ", err)
			//当发现读取数据有错误后，关闭连接
			conn.Close()
			m.waitAuthSocket.Delete(conn.RemoteAddr())
			m.connectedSocket.Delete(conn.GetUid())
			return
		}
		if res != nil {
			fmt.Println(res)
			err = conn.Write(res)
			if err != nil {
				fmt.Printf("write error: %s ", err)
				//switch msg.CmdId {
				//case protocol.AuthReq:
				//	m.retryList.AddRetryMsg(&utils.WaitAckMsg{Msg: msg})
				//case protocol.LogoutReq:
				//case protocol.PullMsgReq:
				//default:
				//
				//}
			}
		}
		conn.SetPingedAt(utils.NowMillisecond())
	}
}

func (m *Handler) ping(conn conn.Conn) (*protocol.NimProtocol, error) {
	return protocol.MakePong(), nil
}
func (m *Handler) authReq(conn conn.Conn, packet *protocol.NimProtocol) (*protocol.NimProtocol, error) {
	msg, err := protocol.Unmarshal(packet)
	res := im.AuthRes{Code: 200, Msg: "验证成功"}
	if err != nil {
		res.Code = 500
		res.Msg = "认证消息体解析异常"
		return protocol.MakeAuthRes(&res), nil
	}
	authReq, _ := msg.(*im.AuthReq)
	if authReq.Token == "" || authReq.Uid == "" {
		res.Code = 500
		res.Msg = "缺少参数"
		return protocol.MakeAuthRes(&res), nil
	}
	conn.SetLogin(true)
	//登陆成功后，将连接转移到已连接队列
	m.waitAuthSocket.Delete(conn.RemoteAddr())
	m.connectedSocket.Store(authReq.Uid, conn)
	//做token跟uid的验证操作，如果验证成功，则将连接设置为已验证的连接，如果验证失败，则返回报错
	return protocol.MakeAuthRes(&res), nil
}
func (m *Handler) logoutReq(conn conn.Conn, packet *protocol.NimProtocol) (*protocol.NimProtocol, error) {
	msg, err := protocol.Unmarshal(packet)
	res := im.LogoutRes{}
	if err != nil {
		res.Code = 500
		res.Msg = "认证消息体解析异常"
		return protocol.MakeLogoutRes(&res), nil
	}
	logoutReq, _ := msg.(*im.LogoutReq)
	//如果消息强转失败，说明服务器有问题
	if logoutReq.Token == "" || logoutReq.Uid == "" {
		res.Code = 500
		res.Msg = "缺少参数"
		return protocol.MakeLogoutRes(&res), nil
	}
	//做登陆退出动作
	conn.Close()
	m.connectedSocket.Delete(conn.GetUid())
	return nil, nil
}
func (m *Handler) pullMsgReq(conn conn.Conn, packet *protocol.NimProtocol) (*protocol.NimProtocol, error) {
	msg, err := protocol.Unmarshal(packet)
	res := im.PullMsgRes{}
	if err != nil {
		res.Code = 500
		res.Msg = "认证消息体解析异常"
		return protocol.MakePullMessageRes(&res), nil
	}
	pullMsgReq, _ := msg.(*im.PullMsgReq)

	if pullMsgReq.Uid == "" || pullMsgReq.Limit == 0 {
		res.Code = 500
		res.Msg = "缺少参数"
		return protocol.MakePullMessageRes(&res), nil
	}
	return nil, nil
}
func (m *Handler) messageSend(conn conn.Conn, packet *protocol.NimProtocol) (ack *protocol.NimProtocol, err error) {
	msg, err := protocol.Unmarshal(packet)
	res := im.MessageSendAck{}
	if err != nil {
		res.Code = 500
		res.Msg = "认证消息体解析异常"
		return protocol.MakeMessageSendAck(&res), nil
	}
	sendMsg, _ := msg.(*im.MessageSend)

	if sendMsg.From == "" || sendMsg.To == "" || sendMsg.TempId == 0 || sendMsg.Type == 0 || sendMsg.Content == "" {
		res.Code = 500
		res.Msg = "缺少参数"
		return protocol.MakeMessageSendAck(&res), nil
	}
	//发送到业务后端，消息即发送完成
	return nil, nil
}
func (m *Handler) notifyAck(conn conn.Conn, packet *protocol.NimProtocol) (*protocol.NimProtocol, error) {
	msg, err := protocol.Unmarshal(packet)
	if err != nil {
		return nil, nil
	}
	notifyAck, _ := msg.(*im.NotifyAck)
	if notifyAck.MsgId == 0 {
		return nil, nil
	}
	//从重试队列里删除消息
	return nil, nil
}
func NewHandler() *Handler {
	handler := &Handler{}
	go handler.loopCheckConnection()
	return handler
}

//循环检查连接，如果发现连接长时间没有活跃，则关闭
func (m *Handler) loopCheckConnection() {
	num := 0
	for {
		//每隔60秒跑一次循环，每次循环中每毫秒检查一个连接，如果以10W个连接算，循环一次检查要100s
		time.Sleep(time.Second * 60)
		m.connectedSocket.Range(func(key, value interface{}) bool {
			socket, _ := value.(conn.Conn)
			//如果此连接的最后一次心跳在5分钟前，则关闭连接
			if socket.GetExpireAt() < utils.NowMillisecond() {
				socket.Close()
				m.connectedSocket.Delete(key)
			}
			num++
			if num > 50 {
				time.Sleep(time.Millisecond * 50)
				num = 0
			}
			return true
		})
		m.waitAuthSocket.Range(func(key, value interface{}) bool {
			socket, _ := value.(conn.Conn)
			//如果此连接的最后一次心跳在5分钟前，则关闭连接
			if socket.GetExpireAt() < utils.NowMillisecond() {
				socket.Close()
				m.waitAuthSocket.Delete(key)
			}
			num++
			if num > 50 {
				time.Sleep(time.Millisecond * 50)
				num = 0
			}
			return true
		})
	}
}
