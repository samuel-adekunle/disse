package beb

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const helloMessageId ds.MessageId = "HelloMessage"
const helloTimerId ds.TimerId = "HelloTimer"

type BebClient struct {
	*ds.AbstractNode
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
	helloTimer := ds.NewTimer(helloTimerId, n.broadcastMessage())
	n.SetTimer(ctx, helloTimer, n.messageDelay)
}

func (n *BebClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Id {
	case helloMessageId:
		fmt.Printf("%v received message: '%v'\n", n.GetAddress(), message.Data.(string))
		return true
	default:
		return false
	}
}

func (n *BebClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Id {
	case helloTimerId:
		n.SendMessage(ctx, timer.Data.(ds.Message), n.bebServer)
		return true
	default:
		return false
	}
}
