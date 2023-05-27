package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// Echo is the type of message used to send an echo message.
	EchoSend ds.MessageType = "EchoSend"
	// EchoDeliver is the type of message used to indicate that an echo message has been delivered.
	EchoDeliver ds.MessageType = "EchoDeliver"
)

// EchoSendData is the data of a send message.
type EchoSendData struct {
	Message ds.Message
}

// EchoDeliverData is the data of a deliver message.
type EchoDeliverData struct {
	Message ds.Message
}

// EchoNode is a node that echoes messages.
type EchoNode struct {
	*ds.AbstractNode
}

// Init is called when the node is initialized by the simulation.
func (n *EchoNode) Init(ctx context.Context) {}

// HandleMessage is called when the node receives a message.
func (n *EchoNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case EchoSend:
		// TODO: 1. decode message data
		// TODO: 2. send echo deliver message back to sender
		return true
	}
	return false
}

// HandleTimer is called when the node receives a timer.
func (n *EchoNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
