package utils

import (
	"github.com/emirpasic/gods/maps/treemap"
	"im/proto/protocol"
	"sync"
)

type WaitAckMsg struct {
	Msg        *protocol.NimProtocol
	RetryTime  int64 //下一次发起重试的时间
	RetryCount int32 //重试次数
	Uid        int64 //发送给谁
}

type RetryList struct {
	sync.RWMutex
	treeMap *treemap.Map
}

func NewRetryList() *RetryList {
	retryList := RetryList{}
	retryList.treeMap = treemap.NewWithIntComparator()
	return &retryList
}
