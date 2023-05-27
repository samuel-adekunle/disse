package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// HelloTimer is the type of timer used to send a hello to an echo node after 1 second.
	HelloTimer ds.TimerType = "HelloTimer"
	// EchoSend is the type of message used to send a hello.
	Hello ds.MessageType = "Hello"
)

// HelloData is the data of a hello message.
type HelloData string

// HelloNode is a node that sends a hello to another node every 1 second.
type HelloNode struct {
	*ds.AbstractNode
	receiver ds.Address
}

// Init is called when the node is initialized by the simulation.
func (n *HelloNode) Init(ctx context.Context) {
	timer := ds.NewTimer(HelloTimer, nil)
	n.SetTimer(ctx, timer, 1*time.Second)
}

// HandleMessage is called when the node receives a message.
func (n *HelloNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case EchoDeliver:
		data := message.Data.(EchoDeliverData)
		fmt.Printf("%s received EchoDeliver: %s\n", n.GetAddress(), data)
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *HelloNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	case HelloTimer:
		echoSendMessage := ds.NewMessage(EchoSend, EchoSendData{
			Message: ds.NewMessage(Hello, HelloData("Hey Jonah")),
		})
		n.SendMessage(ctx, echoSendMessage, n.receiver)
		return true
	default:
		return false
	}
}
