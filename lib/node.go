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
	address  Address
	sim      *Simulation
	subNodes map[Address]Node
}

func NewBaseNode(sim *Simulation, address Address) *BaseNode {
	return &BaseNode{
		address:  address,
		sim:      sim,
		subNodes: make(map[Address]Node),
	}
}

func (n *BaseNode) AddSubNode(address Address, node Node) {
	n.subNodes[address] = node
}

func (n *BaseNode) SubNodesInit(ctx context.Context) {
	var wg sync.WaitGroup
	for address, node := range n.subNodes {
		wg.Add(1)
		go func(_address Address, _node Node) {
			log.Printf("Init(%v)\n", _address)
			_node.Init(ctx)
			_node.SubNodesInit(ctx)
			wg.Done()
		}(address, node)
	}
	wg.Wait()
}

func (n *BaseNode) SubNodesHandleMessage(ctx context.Context, mt MessageTriplet) {
	for address, node := range n.subNodes {
		if address == mt.To {
			log.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message)
			node.HandleMessage(ctx, mt.Message, mt.From)
		} else {
			node.SubNodesHandleMessage(ctx, mt)
		}
	}
}

func (n *BaseNode) SubNodesHandleTimer(ctx context.Context, tt TimerTriplet) {
	for address, node := range n.subNodes {
		if address == tt.From {
			log.Printf("HandleTimer(%v, %v, %v)\n", tt.From, tt.Timer, tt.Duration)
			node.HandleTimer(ctx, tt.Timer, tt.Duration)
		} else {
			node.SubNodesHandleTimer(ctx, tt)
		}
	}
}

func (n *BaseNode) SendMessage(ctx context.Context, message Message, to Address) {
	select {
	case <-ctx.Done():
		log.Printf("StopSim.SendMessage(%v -> %v, %v)\n", n.address, to, message)
		return
	default:
		log.Printf("SendMessage(%v -> %v, %v)\n", n.address, to, message)
		mt := MessageTriplet{message, n.address, to}
		if to == n.address {
			n.sim.HandleMessage(ctx, mt)
		} else if to.Root() == n.address {
			n.SubNodesHandleMessage(ctx, mt)
		} else {
			n.sim.MessageQueue <- mt
		}
	}
}

func (n *BaseNode) SetTimer(ctx context.Context, timer Timer, duration time.Duration) {
	select {
	case <-ctx.Done():
		log.Printf("StopSim.SetTimer(%v, %v, %v)\n", n.address, timer, duration)
		return
	default:
		log.Printf("SetTimer(%v, %v, %v)\n", n.address, timer, duration)
		n.sim.TimerQueue <- TimerTriplet{timer, n.address, duration}
	}
}
