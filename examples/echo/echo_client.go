package main

import (
	"context"
	ds "disse/lib"
	"time"
)

const EchoTimer ds.Timer = "EchoTimer"

type EchoClient struct {
	ds.BaseNode
	EchoServerAddress ds.Address
	EchoInterval      time.Duration
	EchoMessage       ds.Message
	EchoCounter       int
}

func (n *EchoClient) Init(ctx context.Context) {
	n.SetTimer(ctx, EchoTimer, n.EchoInterval)
}

func (n *EchoClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	if message == n.EchoMessage {
		n.EchoCounter++
	}
}

func (n *EchoClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	if timer == EchoTimer {
		n.SendMessage(ctx, n.EchoMessage, n.EchoServerAddress)
		n.SetTimer(ctx, EchoTimer, n.EchoInterval)
	}
}
