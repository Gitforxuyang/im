package utils

import (
	"im/proto/protocol"
	"im/utils/skiplist"
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
	skipList *skiplist.SkipList
}

func NewRetryList() *RetryList {
	retryList := RetryList{}
	list := skiplist.New()
	retryList.skipList = &list
	return &retryList
}
