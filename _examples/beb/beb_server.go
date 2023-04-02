package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const broadcastMessageId ds.MessageType = "BroadcastMessage"

type BebServer struct {
	*ds.AbstractNode
	nodes []ds.Address
	Sent  []ds.Message
}

func (n *BebServer) Init(ctx context.Context) {}

func (n *BebServer) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case broadcastMessageId:
		broadcastMessage := message.Data.(ds.Message)
		for _, node := range n.nodes {
			n.SendMessage(ctx, broadcastMessage, node)
		}
		n.Sent = append(n.Sent, broadcastMessage)
		return true
	default:
		return false
	}
}

func (n *BebServer) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
