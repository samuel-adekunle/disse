package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

// broadcastMessageType is the type of message used to broadcast a message to all nodes.
const broadcastMessageType ds.MessageType = "BebBroadcast"

// deliverMessageType is the type of message used to indicate that a message has been broadcasted.
const deliverMessageType ds.MessageType = "BebDeliver"

// BebMessageData is the data of a broadcast or deliver message.
type BebMessageData struct {
	Message ds.Message
}

// BebNode is a node that broadcasts messages to all other nodes.
type BebNode struct {
	*ds.AbstractNode
	nodes           []ds.Address
	handledMessages map[ds.MessageId]int
}

// Init is called when the node is initialized by the simulation.
func (n *BebNode) Init(ctx context.Context) {
	n.handledMessages = make(map[ds.MessageId]int)
}

// HandleMessage is called when the node receives a message.
//
// If the message is a broadcast message, the message is broadcasted to all other nodes.
//
// All other messages are dropped.
func (n *BebNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case broadcastMessageType:
		data := message.Data.(BebMessageData)
		n.BroadcastMessage(ctx, data.Message, n.nodes)
		deliverMessage := ds.NewMessage(deliverMessageType, data.Message)
		n.SendMessage(ctx, deliverMessage, from)
		return true
	case helloMessageType:
		fmt.Printf("%v received hello: '%v'\n", n.GetAddress(), message.Data.(HelloMessageData))
		n.handledMessages[message.Id]++
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
//
// All timers are dropped.
func (n *BebNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
