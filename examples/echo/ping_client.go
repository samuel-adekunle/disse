package main

import (
	"context"
	ds "disse/lib"
	"time"
)

const PingTimer ds.Timer = "PingTimer"

type PingClient struct {
	ds.BaseNode
	PingMessage   ds.Message
	PongMessage   ds.Message
	ServerAddress ds.Address
	PingInterval  time.Duration
	PongCounter   int
}

func (n *PingClient) Init(ctx context.Context) {
	n.LogInit()
	n.SetTimer(ctx, PingTimer, n.PingInterval)
}

func (n *PingClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	n.BaseNode.LogHandleMessage(message, from)
	if message == n.PongMessage {
		n.PongCounter++
	}
}

func (n *PingClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	n.BaseNode.LogHandleTimer(timer, length)
	if timer == PingTimer {
		n.SendMessage(ctx, n.PingMessage, n.ServerAddress)
		n.SetTimer(ctx, PingTimer, n.PingInterval)
	}
}
