package main

import (
	"context"
	ds "disse/lib"
	"time"
)

type EchoServer struct {
	ds.BaseNode
	*PingServer
	EchoMessage ds.Message
	EchoCounter int
}

func (n *EchoServer) Init(ctx context.Context) {}

func (n *EchoServer) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	n.PingServer.HandleMessage(ctx, message, from)
	if message == n.EchoMessage {
		n.SendMessage(ctx, message, from)
		n.EchoCounter++
	}
}

func (n *EchoServer) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	n.PingServer.HandleTimer(ctx, timer, length)
}
