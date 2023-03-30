package lib

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Node interface {
	// Standard implementations exists (see BaseNode)
	GetState() NodeState
	AddSubNode(Address, Node)
	SubNodesInit(context.Context)
	SubNodesHandleMessage(context.Context, MessageTriplet)
	SubNodesHandleTimer(context.Context, TimerTriplet)
	SubNodesHandleInterrupt(context.Context, InterruptPair) error
	HandleInterrupt(context.Context, Interrupt) error
	SendMessage(ctx context.Context, message Message, to Address)
	BroadcastMessage(ctx context.Context, message Message, to []Address)
	SetTimer(ctx context.Context, timer Timer, duration time.Duration)

	// To be implemented by user for specific node functionality
	Init(context.Context)
	HandleMessage(context.Context, Message, Address)
	HandleTimer(context.Context, Timer, time.Duration)
}

type NodeState string

const (
	Stopped  NodeState = "Stopped"
	Running  NodeState = "Running"
	Sleeping NodeState = "Sleeping"
)

type BaseNode struct {
	address  Address
	sim      *Simulation
	subNodes map[Address]Node
	state    NodeState
}

func NewBaseNode(sim *Simulation, address Address) *BaseNode {
	return &BaseNode{
		address:  address,
		sim:      sim,
		subNodes: make(map[Address]Node),
		state:    Running,
	}
}

func (n *BaseNode) GetState() NodeState {
	return n.state
}

func (n *BaseNode) AddSubNode(address Address, node Node) {
	n.subNodes[address] = node
}

func (n *BaseNode) SubNodesInit(ctx context.Context) {
	var wg sync.WaitGroup
	for address, node := range n.subNodes {
		wg.Add(1)
		go func(_address Address, _node Node) {
			n.sim.debugLog.Printf("Init(%v)\n", _address)
			_node.Init(ctx)
			_node.SubNodesInit(ctx)
			wg.Done()
		}(address, node)
	}
	wg.Wait()
}

func (n *BaseNode) SubNodesHandleMessage(ctx context.Context, mt MessageTriplet) {
	if node, ok := n.subNodes[mt.To]; ok {
		if node.GetState() != Running {
			n.sim.debugLog.Printf("DroppedMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message.Id)
			return
		}
		n.sim.debugLog.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message.Id)
		node.HandleMessage(ctx, mt.Message, mt.From)
	} else {
		var wg sync.WaitGroup
		for _, node := range n.subNodes {
			wg.Add(1)
			go func(_node Node) {
				_node.SubNodesHandleMessage(ctx, mt)
				wg.Done()
			}(node)
		}
		wg.Wait()
	}
}

func (n *BaseNode) SubNodesHandleTimer(ctx context.Context, tt TimerTriplet) {
	if node, ok := n.subNodes[tt.To]; ok {
		if node.GetState() != Running {
			n.sim.debugLog.Printf("DroppedTimer(%v, %v, %v)\n", tt.To, tt.Timer.Id, tt.Duration)
			return
		}
		n.sim.debugLog.Printf("HandleTimer(%v, %v, %v)\n", tt.To, tt.Timer.Id, tt.Duration)
		node.HandleTimer(ctx, tt.Timer, tt.Duration)
	} else {
		var wg sync.WaitGroup
		for _, node := range n.subNodes {
			wg.Add(1)
			go func(_node Node) {
				_node.SubNodesHandleTimer(ctx, tt)
				wg.Done()
			}(node)
		}
		wg.Wait()
	}
}

func (n *BaseNode) SubNodesHandleInterrupt(ctx context.Context, ip InterruptPair) (err error) {
	if node, ok := n.subNodes[ip.To]; ok {
		if node.GetState() == Stopped {
			n.sim.debugLog.Printf("DroppedInterrupt(%v, %v)\n", ip.To, ip.Interrupt.Id)
			return
		}
		n.sim.debugLog.Printf("HandleInterrupt(%v, %v)\n", n.address, ip.Interrupt.Id)
		err = node.HandleInterrupt(ctx, ip.Interrupt)
	} else {
		var wg sync.WaitGroup
		for _, node := range n.subNodes {
			wg.Add(1)
			go func(_node Node) {
				_node.SubNodesHandleInterrupt(ctx, ip)
				wg.Done()
			}(node)
		}
		wg.Wait()
	}
	return
}

func (n *BaseNode) HandleInterrupt(ctx context.Context, interrupt Interrupt) error {
	switch interrupt.Id {
	case StopInterrupt:
		n.state = Stopped
	case SleepInterrupt:
		data := interrupt.Data.(SleepInterruptData)
		n.state = Sleeping
		// Review: Should some sync primitive be used here to make goroutine safe?
		go func() {
			<-time.After(data.Duration)
			startInterrupt := Interrupt{StartInterrupt, nil}
			n.sim.HandleInterrupt(ctx, InterruptPair{startInterrupt, n.address})
		}()
	case StartInterrupt:
		n.state = Running
	default:
		return fmt.Errorf("unknown interrupt: %v", interrupt.Id)
	}
	return nil
}

func (n *BaseNode) SendMessage(ctx context.Context, message Message, to Address) {
	select {
	case <-ctx.Done():
		n.sim.debugLog.Printf("StopSim.SendMessage(%v -> %v, %v)\n", n.address, to, message.Id)
		return
	default:
		n.sim.umlLog.Printf("%v -> %v : %v\n", n.address, to, message.Id)
		n.sim.debugLog.Printf("SendMessage(%v -> %v, %v)\n", n.address, to, message.Id)
		mt := MessageTriplet{message, n.address, to}
		if to.Root() == n.address.Root() {
			n.sim.HandleMessage(ctx, mt)
		} else {
			n.sim.MessageQueue <- mt
		}
	}
}

func (n *BaseNode) BroadcastMessage(ctx context.Context, message Message, to []Address) {
	select {
	case <-ctx.Done():
		n.sim.debugLog.Printf("StopSim.BroadcastMessage(%v -> %v, %v)\n", n.address, to, message.Id)
		return
	default:
		n.sim.debugLog.Printf("BroadcastMessage(%v -> %v, %v)\n", n.address, to, message.Id)
		for _, address := range to {
			n.SendMessage(ctx, message, address)
		}
	}
}

func (n *BaseNode) SetTimer(ctx context.Context, timer Timer, duration time.Duration) {
	select {
	case <-ctx.Done():
		n.sim.debugLog.Printf("StopSim.SetTimer(%v, %v, %v)\n", n.address, timer.Id, duration)
		return
	default:
		n.sim.umlLog.Printf("%v -> %v : %v\n", n.address, n.address, timer.Id)
		n.sim.debugLog.Printf("SetTimer(%v, %v, %v)\n", n.address, timer.Id, duration)
		n.sim.TimerQueue <- TimerTriplet{timer, n.address, duration}
	}
}
