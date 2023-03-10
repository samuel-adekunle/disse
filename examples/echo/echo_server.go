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

func (n *EchoServer) Init(ctx context.Context) {
	n.LogInit()
	n.PingServer.Init(ctx)
}

func (n *EchoServer) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	n.BaseNode.LogHandleMessage(message, from)
	n.PingServer.HandleMessage(ctx, message, from)
	if message == n.EchoMessage {
		n.SendMessage(ctx, message, from)
		n.EchoCounter++
	}
}

func (n *EchoServer) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	n.BaseNode.LogHandleTimer(timer, length)
	n.PingServer.HandleTimer(ctx, timer, length)
}
