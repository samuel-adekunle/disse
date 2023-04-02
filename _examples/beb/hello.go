package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

// helloMessageType is the type of message used to send a hello message.
const helloMessageType ds.MessageType = "HelloMessage"

// HelloMessageData is the data of a hello message.
type HelloMessageData string

// helloTimerType is the type of timer used to send a hello message after an initial delay.
const helloTimerType ds.TimerType = "HelloTimer"

// HelloTimerData is the data of a hello timer.
type HelloTimerData ds.Message

// HelloNode is a node that sends a hello message to all nodes in the network after an initial delay.
//
// The hello message is sent using the Best-Effort Broadcast (BEB) primitive.
type HelloNode struct {
	*ds.AbstractNode
	beb             ds.Address
	sendAfter       time.Duration
	handledMessages map[ds.MessageId]int
	sentMessages    []ds.MessageId
}

// Init is called when the node is initialized by the simulation.
//
// A timer is set to send a hello message after an initial delay.
func (n *HelloNode) Init(ctx context.Context) {
	n.handledMessages = make(map[ds.MessageId]int)
	n.sentMessages = make([]ds.MessageId, 0)
	helloMessage := ds.NewMessage(helloMessageType, HelloMessageData(fmt.Sprintf("Hello from %v", n.GetAddress())))
	broadcastMessage := ds.NewMessage(broadcastMessageType, BebMessageData{Message: helloMessage})
	helloTimer := ds.NewTimer(helloTimerType, HelloTimerData(broadcastMessage))
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
		fmt.Printf("%v received hello: '%v'\n", n.GetAddress(), message.Data.(HelloMessageData))
		n.handledMessages[message.Id]++
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
		broadcastMessage := timer.Data.(HelloTimerData)
		helloMessage := broadcastMessage.Data.(BebMessageData).Message
		n.SendMessage(ctx, ds.Message(broadcastMessage), n.beb)
		n.sentMessages = append(n.sentMessages, helloMessage.Id)
		return true
	default:
		return false
	}
}

// HelloNodeFaultyType is the type of node that is faulty and crashes after a certain time.
const HelloFaultTimerType ds.TimerType = "HelloFaultTimer"

type HelloFaultTimerData ds.Interrupt

// FaultyHelloNode is a HelloNode that is faulty can crashes after a certain time.
type FaultyHelloNode struct {
	*HelloNode
	faultAfter time.Duration
}

// Init is called when the node is initialized by the simulation.
//
// A timer is set to interrupt the node after an initial delay.
func (n *FaultyHelloNode) Init(ctx context.Context) {
	n.HelloNode.Init(ctx)
	interrupt := ds.NewInterrupt(ds.StopInterrupt, nil)
	faultTimer := ds.NewTimer(HelloFaultTimerType, HelloFaultTimerData(interrupt))
	n.SetTimer(ctx, faultTimer, n.faultAfter)
}

// HandleTimer is called when the node receives a timer.
//
// The fault timer is handled by interrupting the node.
func (n *FaultyHelloNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	case HelloFaultTimerType:
		data := timer.Data.(HelloFaultTimerData)
		n.SendInterrupt(ctx, ds.Interrupt(data), n.GetAddress())
		return true
	default:
		return n.HelloNode.HandleTimer(ctx, timer, length)
	}
}
