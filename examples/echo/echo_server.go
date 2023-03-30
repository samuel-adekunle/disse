package echo

import (
	"context"
	ds "disse/lib"
	"time"
)

type EchoServer struct {
	*ds.BaseNode
	echoMessage ds.Message
	EchoCounter int
}

func (n *EchoServer) Init(ctx context.Context) {}

func (n *EchoServer) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	if message.Id == n.echoMessage.Id {
		n.SendMessage(ctx, message, from)
		n.EchoCounter++
	}
}

func (n *EchoServer) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {}
