package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// PingMessageType is the type of message used to send a ping message.
	PingMessageType = "Ping"
	// PongMessageType is the type of message used to send a pong message.
	PongMessageType = "Pong"
)

// PingNode is a node that pings other nodes.
type PingNode struct {
	*ds.LocalNode
	Nodes []ds.Address
}

// Init is called when the node is initialized by the simulation.
func (n *PingNode) Init(ctx context.Context) {
	n.BroadcastMessage(ctx, ds.NewMessage("ping", nil), n.Nodes)
}

// HandleMessage is called when the node receives a message.
func (n *PingNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case PingMessageType:
		n.SendMessage(ctx, ds.NewMessage(PongMessageType, nil), from)
		return true
	case PongMessageType:
		fmt.Printf("%s: Received pong from %s\n", n.GetAddress(), from)
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *PingNode) HandleTimer(ctx context.Context, timer ds.Timer, duration time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
