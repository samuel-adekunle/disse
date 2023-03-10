package main

import (
	"context"
	ds "disse/lib"
	"time"
)

type PingServer struct {
	*ds.BaseNode
	pingMessage ds.Message
	pongMessage ds.Message
	PingCounter int
}

func (n *PingServer) Init(ctx context.Context) {}

func (n *PingServer) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	if message == n.pingMessage {
		n.PingCounter++
		n.SendMessage(ctx, n.pongMessage, from)
	}
}

func (n *PingServer) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {}
