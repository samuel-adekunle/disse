package lib

import (
	"context"
	"log"
	"time"
)

type Node interface {
	Init(context.Context)
	SubNodesInit(context.Context)
	AddSubNode(Address, Node)
	HandleMessage(context.Context, Message, Address)
	HandleTimer(context.Context, Timer, time.Duration)
}

type BaseNode struct {
	Address      Address
	MessageQueue chan MessageTriplet
	TimerQueue   chan TimerTriplet
	SubNodes     map[Address]Node
}

func NewBaseNode(sim *Simulation, address Address) BaseNode {
	return BaseNode{
		Address:      address,
		MessageQueue: sim.MessageQueue,
		TimerQueue:   sim.TimerQueue,
		SubNodes:     make(map[Address]Node),
	}
}

func (n *BaseNode) SubNodesInit(ctx context.Context) {
	for address, node := range n.SubNodes {
		log.Printf("Init(%v)\n", address)
		node.Init(ctx)
	}
}

func (n *BaseNode) AddSubNode(address Address, node Node) {
	n.SubNodes[address] = node
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
