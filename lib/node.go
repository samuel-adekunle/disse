package lib

import (
	"log"
	"time"
)

// Node type, for all nodes in the system
type Node interface {
	Init()
	HandleMessage(Message, Address)
	HandleTimer(Timer)
}

// BaseNode
type BaseNode struct {
	Address      Address
	MessageQueue chan MessageTriplet
	TimerQueue   chan TimerTriplet
}

func (n *BaseNode) Init() {
	panic("implement me")
}

func (n *BaseNode) HandleMessage(message Message, from Address) {
	panic("implement me")
}

func (n *BaseNode) HandleTimer(timer Timer) {
	panic("implement me")
}

func (n *BaseNode) SendMessage(message Message, to Address) {
	log.Printf("SendMessage(%v -> %v, %v)\n", n.Address, to, message)
	n.MessageQueue <- MessageTriplet{message, n.Address, to}
}

func (n *BaseNode) SetTimer(timer Timer, length time.Duration) {
	log.Printf("SetTimer(%v, %v, %v)\n", n.Address, timer, length)
	n.TimerQueue <- TimerTriplet{timer, n.Address, length}
}
