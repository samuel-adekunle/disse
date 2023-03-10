package main

import (
	"context"
	ds "disse/lib"
	"time"
)

type EchoClient struct {
	ds.BaseNode
	*PingClient
}

func (n *EchoClient) Init(ctx context.Context) {
	n.LogInit()
	n.PingClient.Init(ctx)
}

func (n *EchoClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	n.BaseNode.LogHandleMessage(message, from)
	n.PingClient.HandleMessage(ctx, message, from)
}

func (n *EchoClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	n.BaseNode.LogHandleTimer(timer, length)
	n.PingClient.HandleTimer(ctx, timer, length)
}
