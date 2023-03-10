package main

import (
	"context"
	ds "disse/lib"
	"time"
)

const EchoTimer ds.Timer = "EchoTimer"

type EchoClient struct {
	ds.BaseNode
	*PingClient
	EchoServerAddress ds.Address
	EchoInterval      time.Duration
	EchoMessage       ds.Message
	EchoCounter       int
}

func (n *EchoClient) Init(ctx context.Context) {
	n.LogInit()
	n.PingClient.Init(ctx)
	n.SetTimer(ctx, EchoTimer, n.EchoInterval)
}

func (n *EchoClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	n.BaseNode.LogHandleMessage(message, from)
	n.PingClient.HandleMessage(ctx, message, from)
	if message == n.EchoMessage {
		n.EchoCounter++
	}
}

func (n *EchoClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	n.BaseNode.LogHandleTimer(timer, length)
	n.PingClient.HandleTimer(ctx, timer, length)
	if timer == EchoTimer {
		n.SendMessage(ctx, n.EchoMessage, n.EchoServerAddress)
		n.SetTimer(ctx, EchoTimer, n.EchoInterval)
	}
}
