package echo

import (
	"context"
	ds "disse/lib"
	"time"
)

var pingTimer ds.Timer = ds.NewTimer(ds.TimerId("PingTimer"), nil)

type PingClient struct {
	*ds.BaseNode
	pingMessage   ds.Message
	pongMessage   ds.Message
	serverAddress ds.Address
	pingInterval  time.Duration
	PongCounter   int
}

func (n *PingClient) Init(ctx context.Context) {
	n.SetTimer(ctx, pingTimer, n.pingInterval)
}

func (n *PingClient) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	if message.Id == n.pongMessage.Id {
		n.PongCounter++
	}
}

func (n *PingClient) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {
	if timer.Id == pingTimer.Id {
		n.SendMessage(ctx, n.pingMessage, n.serverAddress)
		n.SetTimer(ctx, pingTimer, n.pingInterval)
	}
}
