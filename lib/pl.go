package lib

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// PlSend is the type of message used to send a message to another node.
	PlSend ds.MessageType = "PlSend"
	// PlDeliver is the type of message used to indicate that a message has been sent.
	PlDeliver ds.MessageType = "PlDeliver"
)

// PlSendData is the data of a send message.
type PlSendData struct {
	Destination ds.Address
	Message     ds.Message
}

// PlDeliverData is the data of a deliver message.
type PlDeliverData struct {
	Source  ds.Address
	Message ds.Message
}

// PlNode is a node that implements perfect point-to-point links.
//
// This implementation uses the "Eliminate Duplicates" algorithm.
type PlNode struct {
	*ds.LocalNode
	deliveredMessages map[ds.MessageId]bool
}

// Init is called when the node is initialized by the simulation.
func (n *PlNode) Init(ctx context.Context) {
	n.deliveredMessages = make(map[ds.MessageId]bool)
}

// HandleMessage is called when the node receives a message.
func (n *PlNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case PlSend:
		data := message.Data.(PlSendData)
		if _, ok := n.deliveredMessages[message.Id]; ok {
			return true
		}
		deliverMessage := ds.NewMessage(PlDeliver, PlDeliverData{
			Source:  from,
			Message: data.Message,
		})
		n.SendMessage(ctx, deliverMessage, from)
		n.deliveredMessages[message.Id] = true
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *PlNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
