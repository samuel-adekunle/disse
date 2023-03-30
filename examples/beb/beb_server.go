package beb

import (
	"context"
	ds "disse/lib"
	"time"
)

const broadcastMessageId ds.MessageId = "BroadcastMessage"

type BebServer struct {
	*ds.BaseNode
	nodes []ds.Address
}

func (n *BebServer) Init(ctx context.Context) {}

func (n *BebServer) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {
	if message.Id == broadcastMessageId {
		data := message.Data.(ds.Message)
		for _, node := range n.nodes {
			n.SendMessage(ctx, data, node)
		}
	}
}

func (n *BebServer) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {}
