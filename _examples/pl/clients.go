package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const helloMessageType ds.MessageType = "HelloMsg"

type SenderNode struct {
	*ds.AbstractNode
	pl       ds.Address
	receiver ds.Address
}

func (n *SenderNode) Init(ctx context.Context) {
	helloMessage := ds.NewMessage(helloMessageType, fmt.Sprintf("Hello from %v", n.GetAddress()))
	sendMessage := ds.NewMessage(sendMessageType, PLMessageData{
		Message: helloMessage,
		To:      n.receiver,
	})
	n.SendMessage(ctx, sendMessage, n.pl)
}

func (n *SenderNode) HandleMessage(ctx context.Context, message ds.Message, address ds.Address) bool {

}

func (n *SenderNode) HandleTimer(ctx context.Context, timer Timer, duration time.Duration) bool {

}
