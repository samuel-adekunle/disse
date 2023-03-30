package beb

import (
	"context"
	ds "disse"
	"fmt"
	"time"
)

const helloMessageId ds.MessageId = "HelloMessage"
const helloTimerId ds.TimerId = "HelloTimer"

type BebClient struct {
	*ds.BaseNode
	bebServer    ds.Address
	messageDelay time.Duration
}

func (n *BebClient) broadcastMessage() ds.Message {
	return ds.NewMessage(
		broadcastMessageId,
		ds.NewMessage(helloMessageId, fmt.Sprintf("Hello from %s", n.GetAddress())),
	)
}

func (n *BebClient) Init(ctx context.Context) {
	var helloTimer ds.Timer = ds.NewTimer(helloTimerId, n.broadcastMessage())
	n.SetTimer(ctx, helloTimer, n.messageDelay)
}

func (n *BebClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	if message.Id == helloMessageId {
		fmt.Printf("%v received message: '%v'\n", n.GetAddress(), message.Data.(string))
	}
}

func (n *BebClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	if timer.Id == helloTimerId {
		n.SendMessage(ctx, timer.Data.(ds.Message), n.bebServer)
	}
}
