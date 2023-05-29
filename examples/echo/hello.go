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

// HelloNode is a node that sends a hello to an echo node after 1 second.
type HelloNode struct {
	*ds.LocalNode
	echoNode ds.Address
}

func (n *HelloNode) setHelloTimer(ctx context.Context, duration time.Duration) {
	timer := ds.NewTimer(HelloTimer, nil)
	n.SetTimer(ctx, timer, duration)
}

// Init is called when the node is initialized by the simulation.
func (n *HelloNode) Init(ctx context.Context) {
	n.setHelloTimer(ctx, 1*time.Second)
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
			Message: ds.NewMessage(Hello, HelloData("Hello DISSE!")),
		})
		n.SendMessage(ctx, echoSendMessage, n.echoNode)
		n.setHelloTimer(ctx, 1*time.Second)
		return true
	default:
		return false
	}
}
