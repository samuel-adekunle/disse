package lib

import (
	"context"
	"log"
	"sync"
	"time"
)

type Node interface {
	Init(context.Context)
	SubNodesInit(context.Context)
	AddSubNode(Address, Node)
	HandleMessage(context.Context, Message, Address)
	HandleTimer(context.Context, Timer, time.Duration)
	SubNodesHandleMessage(context.Context, MessageTriplet)
	SubNodesHandleTimer(context.Context, TimerTriplet)
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

func (n *BaseNode) AddSubNode(address Address, node Node) {
	n.SubNodes[address] = node
}

func (n *BaseNode) SubNodesInit(ctx context.Context) {
	for address, node := range n.SubNodes {
		log.Printf("Init(%v)\n", address)
		node.Init(ctx)
	}
}

func (n *BaseNode) SubNodesHandleMessage(ctx context.Context, mt MessageTriplet) {
	var wg sync.WaitGroup
	for address, node := range n.SubNodes {
		wg.Add(1)
		go func(_address Address, _node Node) {
			if _address == mt.To {
				log.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message)
				_node.HandleMessage(ctx, mt.Message, mt.From)
			} else {
				_node.SubNodesHandleMessage(ctx, mt)
			}
		}(address, node)
	}
	wg.Wait()
}

func (n *BaseNode) SubNodesHandleTimer(ctx context.Context, tt TimerTriplet) {
	var wg sync.WaitGroup
	for address, node := range n.SubNodes {
		wg.Add(1)
		go func(_address Address, _node Node) {
			if _address == tt.From {
				log.Printf("HandleTimer(%v, %v, %v)\n", tt.From, tt.Timer, tt.Length)
				_node.HandleTimer(ctx, tt.Timer, tt.Length)
			} else {
				_node.SubNodesHandleTimer(ctx, tt)
			}
		}(address, node)
	}
	wg.Wait()
}

func (n *BaseNode) SendMessage(ctx context.Context, message Message, to Address) {
	select {
	case <-ctx.Done():
		log.Printf("StopSim.SendMessage(%v -> %v, %v)\n", n.Address, to, message)
		return
	default:
		log.Printf("SendMessage(%v -> %v, %v)\n", n.Address, to, message)
		n.MessageQueue <- MessageTriplet{message, n.Address, to}
	}
}

func (n *BaseNode) SetTimer(ctx context.Context, timer Timer, length time.Duration) {
	select {
	case <-ctx.Done():
		log.Printf("StopSim.SetTimer(%v, %v, %v)\n", n.Address, timer, length)
		return
	default:
		log.Printf("SetTimer(%v, %v, %v)\n", n.Address, timer, length)
		n.TimerQueue <- TimerTriplet{timer, n.Address, length}
	}
}
