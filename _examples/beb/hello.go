package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

// helloMessageType is the type of message used to send a hello message.
const helloMessageType ds.MessageType = "HelloMessage"

// helloTimerType is the type of timer used to send a hello message after an initial delay.
const helloTimerType ds.TimerType = "HelloTimer"

// HelloNode is a node that sends a hello message to all nodes in the network after an initial delay.
//
// The hello message is sent using the Best-Effort Broadcast (BEB) primitive.
type HelloNode struct {
	*ds.AbstractNode
	beb       ds.Address
	sendAfter time.Duration
}

// Init is called when the node is initialized by the simulation.
//
// A timer is set to send a hello message after an initial delay.
func (n *HelloNode) Init(ctx context.Context) {
	helloMessage := ds.NewMessage(helloMessageType, fmt.Sprintf("Hello from %v", n.GetAddress()))
	broadcastMessage := ds.NewMessage(broadcastMessageType, BebMessageData{Message: helloMessage})
	helloTimer := ds.NewTimer(helloTimerType, broadcastMessage)
	n.SetTimer(ctx, helloTimer, n.sendAfter)
}

// HandleMessage is called when the node receives a message.
//
// If the message is a hello message, the message is printed.
//
// All other messages are dropped.
func (n *HelloNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case helloMessageType:
		fmt.Printf("%v received message: '%v'\n", n.GetAddress(), message.Data.(string))
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
//
// If the timer is a hello timer, the timer is handled.
//
// All other timers are dropped.
func (n *HelloNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	case helloTimerType:
		n.SendMessage(ctx, timer.Data.(ds.Message), n.beb)
		return true
	default:
		return false
	}
}
