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
	Message ds.Message
}

// PlNode is a node that implements perfect point-to-point links.
//
// Note that this node does not need to be used as the internal `SendMessage` method
// is equivalent to a perfect point-to-point link.
//
// This node is provided as an example of how to implement a module using the DISSE library
// given a specification and is not intended to be used in production.
type PlNode struct {
	*ds.AbstractNode
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
		n.SendMessage(ctx, data.Message, data.Destination)
		// XXX(samuel-adekunle): assume that the message is always delivered
		n.deliveredMessages[message.Id] = true
		deliverMessage := ds.NewMessage(PlDeliver, PlDeliverData{message})
		n.SendMessage(ctx, deliverMessage, from)
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
