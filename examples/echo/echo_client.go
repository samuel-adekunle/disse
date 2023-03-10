package echo

import (
	"context"
	ds "disse/lib"
	"time"
)

const echoTimer ds.Timer = "EchoTimer"

type EchoClient struct {
	*ds.BaseNode
	echoServerAddress ds.Address
	echoInterval      time.Duration
	echoMessage       ds.Message
	EchoCounter       int
}

func (n *EchoClient) Init(ctx context.Context) {
	n.SetTimer(ctx, echoTimer, n.echoInterval)
}

func (n *EchoClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	if message == n.echoMessage {
		n.EchoCounter++
	}
}

func (n *EchoClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	if timer == echoTimer {
		n.SendMessage(ctx, n.echoMessage, n.echoServerAddress)
		n.SetTimer(ctx, echoTimer, n.echoInterval)
	}
}
