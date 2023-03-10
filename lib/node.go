package lib

import (
	"context"
	"log"
	"time"
)

type Node interface {
	Init(context.Context)
	HandleMessage(context.Context, Message, Address)
	HandleTimer(context.Context, Timer, time.Duration)
}

type BaseNode struct {
	Address      Address
	MessageQueue chan MessageTriplet
	TimerQueue   chan TimerTriplet
}

func (n *BaseNode) LogInit() {
	log.Printf("Init(%v)\n", n.Address)
}

func (n *BaseNode) LogHandleMessage(message Message, from Address) {
	log.Printf("HandleMessage(%v -> %v, %v)\n", from, n.Address, message)
}

func (n *BaseNode) LogHandleTimer(timer Timer, length time.Duration) {
	log.Printf("HandleTimer(%v, %v, %v)\n", n.Address, timer, length)
}

func (n *BaseNode) SendMessage(ctx context.Context, message Message, to Address) {
	select {
	case <-ctx.Done():
		log.Printf("Timeout.SendMessage(%v -> %v, %v)\n", n.Address, to, message)
		return
	default:
		log.Printf("SendMessage(%v -> %v, %v)\n", n.Address, to, message)
		n.MessageQueue <- MessageTriplet{message, n.Address, to}
	}
}

func (n *BaseNode) SetTimer(ctx context.Context, timer Timer, length time.Duration) {
	select {
	case <-ctx.Done():
		log.Printf("Timeout.SetTimer(%v, %v, %v)\n", n.Address, timer, length)
		return
	default:
		log.Printf("SetTimer(%v, %v, %v)\n", n.Address, timer, length)
		n.TimerQueue <- TimerTriplet{timer, n.Address, length}
	}
}
