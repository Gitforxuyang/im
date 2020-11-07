package utils

import (
	"fmt"
	"im/proto/protocol"
	"im/utils/skiplist"
	"strconv"
	"sync"
)

type WaitAckMsg struct {
	Msg        *protocol.NimProtocol
	RetryTime  int64 //下一次发起重试的时间
	RetryCount int32 //重试次数
	Uid        int64 //发送给谁
}

func (m *WaitAckMsg) ExtractKey() float64 {
	key, _ := strconv.ParseFloat(fmt.Sprintf("%d.%d", m.RetryTime, m.Uid), 64)
	return key
}

func (m *WaitAckMsg) String() string {
	return fmt.Sprintf("%d.%d", m.RetryTime, m.Uid)
}

type RetryList struct {
	sync.RWMutex
	skipList *skiplist.SkipList
}

func NewRetryList() *RetryList {
	retryList := RetryList{}
	list := skiplist.New()
	retryList.skipList = &list
	return &retryList
}

func (m *RetryList) AddRetryMsg(msg *WaitAckMsg) error {
	m.Lock()
	defer m.Unlock()
	m.skipList.Insert(msg)
	return nil
}

func (m *RetryList) RemoveRetryMsg(msg *WaitAckMsg) error {
	m.Lock()
	defer m.Unlock()
	m.skipList.Delete(msg)
	return nil
}

func (m *RetryList) GetWaitRetryMsg() (list []*WaitAckMsg, err error) {
	now := WaitAckMsg{RetryTime: NowMillisecond()}
	mid, ok := m.skipList.FindGreaterOrEqual(&now)
	if ok {
		pre := m.skipList.Prev(mid)
		//如果已经到了尾部了，则证明循环结束了
		if pre == m.skipList.GetLargestNode() {
			len := len(list)
			newList := make([]*WaitAckMsg, len)
			for k, v := range list {
				newList[len-k-1] = v
			}
			return newList, nil
		}
		msg, _ := pre.GetValue().(*WaitAckMsg)
		list = append(list, msg)
	}
	return list, nil
}
