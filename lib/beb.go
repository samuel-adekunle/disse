package lib

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// BebBroadcast is the type of message used to broadcast a message to all nodes.
	BebBroadcast ds.MessageType = "BebBroadcast"
	// BebDeliver is the type of message used to indicate that a message has been broadcast.
	BebDeliver ds.MessageType = "BebDeliver"
)

// BebBroadcastData is the data of a broadcast message.
type BebBroadcastData struct {
	Message ds.Message
}

// BebDeliverData is the data of a deliver message.
type BebDeliverData struct {
	Source  ds.Address
	Message ds.Message
}

// BebNode is a node that implements best-effort broadcast.
//
// This implementation uses the "Basic Broadcast" algorithm and makes no
// assumption on failure detection and message reliability.
type BebNode struct {
	*ds.AbstractNode
	Nodes []ds.Address
}

// Init is called when the node is initialized by the simulation.
func (n *BebNode) Init(ctx context.Context) {}

// HandleMessage is called when the node receives a message.
func (n *BebNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case BebBroadcast:
		data := message.Data.(BebBroadcastData)
		deliverMessage := ds.NewMessage(BebDeliver, BebDeliverData{
			Source:  from,
			Message: data.Message,
		})
		// XXX(samuel-adekunle): assumes that a message is always delivered.
		// correct implementation uses perfect point-to-point links.
		n.BroadcastMessage(ctx, deliverMessage, n.Nodes)
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *BebNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
