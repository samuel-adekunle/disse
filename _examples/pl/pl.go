package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	sendTimerType      ds.TimerType   = "PLSendTimer"
	sendMessageType    ds.MessageType = "PLSendMessage"
	deliverMessageType ds.MessageType = "PLDeliverMessage"
)

type PLMessageData struct {
	Message ds.Message
	From    ds.Address
	To      ds.Address
}

type PerfectLink struct {
	*ds.AbstractNode
	sendInterval time.Duration
	messages     []PLMessageData
}

func (n *PerfectLink) Init(ctx context.Context) {
	n.messages = make([]PLMessageData, 0)
	sendTimer := ds.NewTimer(sendTimerType, nil)
	n.SetTimer(ctx, sendTimer, n.sendInterval)
}

func (n *PerfectLink) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case sendMessageType:
		data := message.Data.(PLMessageData)
		n.SendMessage(ctx, data.Message, data.To)
		n.messages = append(n.messages, data)
		return true
	case deliverMessageType:
		data := message.Data.(PLMessageData)
		n.SendMessage(ctx, message, data.From)
		// TODO: filter data from sentmessages
		return true
	default:
		return false
	}
}

func (n *PerfectLink) HandleTimer(ctx context.Context, timer ds.Timer) bool {
	switch timer.Type {
	case sendTimerType:
		for _, data := range n.messages {
			n.SendMessage(ctx)
		}
		return true
	default:
		return false
	}
}
