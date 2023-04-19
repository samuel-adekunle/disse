package disse

import (
	"context"
	"sync"
	"time"
)

// Node is the interface that must be implemented by all nodes in the distributed system.
//
// Most of the methods are implemented by AbstractNode and should not be overridden unless you know what you are doing.
//
// The Init, HandleMessage and HandleTimer methods should be implemented by concrete nodes to provide the desired behaviour.
type Node interface {
	// Implemented by AbstractNode:
	GetAddress() Address
	GetState() NodeState
	GetSimulation() *Simulation

	GetSubNodes() map[Address]Node
	GetSubNode(Address) Node
	AddSubNode(Address, Node)
	RemoveSubNode(Address)

	InitAll(context.Context)
	FindMessageHandler(context.Context, MessageTriplet) (handled bool)
	FindTimerHandler(context.Context, TimerTriplet) (handled bool)
	FindInterruptHandler(context.Context, InterruptTriplet) (handled bool)

	SendMessage(context.Context, Message, Address)
	BroadcastMessage(context.Context, Message, []Address)
	SetTimer(context.Context, Timer, time.Duration)
	SendInterrupt(context.Context, Interrupt, Address)

	HandleInterrupt(context.Context, Interrupt, Address) (handled bool)

	// To be implemented by concrete node:

	Init(context.Context)
	HandleMessage(context.Context, Message, Address) (handled bool)
	HandleTimer(context.Context, Timer, time.Duration) (handled bool)
}

// NodeState is a string that represents the state of a node.
//
// It can be either Stopped, Running or Sleeping.
type NodeState int

const (
	// Running is the state of a node that is running.
	Running NodeState = iota
	// Sleeping is the state of a node that is sleeping.
	Sleeping
	// Stopped is the state of a node that is stopped.
	Stopped
)

// AbstractNode is a base implementation of the Node interface.
//
// It provides the basic functionality for sending messages, setting timers and handling interrupts.
//
// It also provides the functionality for handling sub nodes.
//
// It is recommended that all nodes extend this struct to avoid having to implement the same functionality multiple times.
//
// It is not recommended to override any of the methods in this struct unless you know what you are doing.
type AbstractNode struct {
	address  Address
	sim      *Simulation
	subNodes map[Address]Node
	state    NodeState
}

// NewAbstractNode creates a new AbstractNode.
func NewAbstractNode(sim *Simulation, address Address) *AbstractNode {
	return &AbstractNode{
		address:  address,
		sim:      sim,
		subNodes: make(map[Address]Node),
		state:    Running,
	}
}

// GetAddress returns the address of the node.
func (n *AbstractNode) GetAddress() Address {
	return n.address
}

// GetState returns the state of the node.
func (n *AbstractNode) GetState() NodeState {
	return n.state
}

// GetSimulation returns the simulation that the node is part of.
func (n *AbstractNode) GetSimulation() *Simulation {
	return n.sim
}

// GetSubNodes returns a map of all sub nodes.
func (n *AbstractNode) GetSubNodes() map[Address]Node {
	return n.subNodes
}

// GetSubNode returns a sub node with the given address.
func (n *AbstractNode) GetSubNode(address Address) Node {
	return n.subNodes[address]
}

// AddSubNode adds a sub node to the node.
func (n *AbstractNode) AddSubNode(address Address, node Node) {
	n.subNodes[address] = node
}

// RemoveSubNode removes a sub node with the given address.
func (n *AbstractNode) RemoveSubNode(address Address) {
	delete(n.subNodes, address)
}

// InitAll initializes all sub nodes of a node.
func (n *AbstractNode) InitAll(ctx context.Context) {
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			n.sim.initNode(ctx, _node)
		}(node)
	}
	wg.Wait()
}

// FindMessageHandler finds the correct sub node to handle a message.
//
// If the message is for the current sub node, the HandleMessage method is called.
//
// If the message is not for the current sub node, the FindMessageHandler method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the message is dropped.
//
// If a match is found, and the message is handled successfully, the method returns true, otherwise it returns false.
func (n *AbstractNode) FindMessageHandler(ctx context.Context, mt MessageTriplet) (handled bool) {
	if node, ok := n.subNodes[mt.To]; ok {
		return n.sim.handleMessage(ctx, node, mt)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.FindMessageHandler(ctx, mt)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
}

// FindTimerHandler finds the correct sub node to handle a timer.
//
// If the timer is for the current sub node, the HandleTimer method is called.
//
// If the timer is not for the current sub node, the FindTimerHandler method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the timer is dropped.
//
// If a match is found, and the timer is handled successfully, the method returns true, otherwise it returns false.
func (n *AbstractNode) FindTimerHandler(ctx context.Context, tt TimerTriplet) (handled bool) {
	if node, ok := n.subNodes[tt.To]; ok {
		return n.sim.handleTimer(ctx, node, tt)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.FindTimerHandler(ctx, tt)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
}

// FindInterruptHandler handles an interrupt for all sub nodes.
//
// If the interrupt is for the node, the HandleInterrupt method is called.
//
// If the interrupt is not for the node, the FindInterruptHandler method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the interrupt is dropped.
//
// If an unknown interrupt is received, the interrupt is dropped and the function returns false, otherwise true.
func (n *AbstractNode) FindInterruptHandler(ctx context.Context, it InterruptTriplet) (handled bool) {
	if node, ok := n.subNodes[it.To]; ok {
		return n.sim.handleInterrupt(ctx, node, it)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.FindInterruptHandler(ctx, it)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
}

// SendMessage sends a message to another node in the simulation.
//
// If the message is for the node or one of it's sub nodes, the HandleMessage method is called immediately.
//
// If the message is for a node outside the node's root node, the message is sent to the message queue which is handled by the simulation.
//
// The simulation will handle the message after introducing a random amount of latency.
func (n *AbstractNode) SendMessage(ctx context.Context, message Message, to Address) {
	select {
	case <-ctx.Done():
		return
	default:
		n.sim.LogSendMessage(n.address, to, message)
		mt := MessageTriplet{message, n.address, to}
		if to.Root() == n.address.Root() {
			if handled := n.sim.HandleMessage(ctx, mt); !handled {
				n.sim.DropMessage(ctx, mt)
			}
		} else {
			n.sim.messageQueue <- mt
		}
	}
}

// BroadcastMessage sends a message to multiple nodes in the simulation.
//
// See SendMessage for more details on how the message is sent.
func (n *AbstractNode) BroadcastMessage(ctx context.Context, message Message, to []Address) {
	select {
	case <-ctx.Done():
		return
	default:
		for _, address := range to {
			n.SendMessage(ctx, message, address)
		}
	}
}

// SetTimer sets a timer for the node.
//
// The timer is sent to the timer queue which is handled by the simulation.
//
// The simulation will handle the timer after the specified duration and call the HandleTimer method.
func (n *AbstractNode) SetTimer(ctx context.Context, timer Timer, duration time.Duration) {
	select {
	case <-ctx.Done():
		return
	default:
		n.sim.LogSetTimer(n.address, timer, duration)
		n.sim.timerQueue <- TimerTriplet{timer, n.address, duration}
	}
}

// SendInterrupt sends an interrupt to another node in the simulation.
//
// Interrupts are handled immediately and do not go through the message queue.
//
// See HandleInterrupt for more details on how interrupts are handled.
//
// If the interrupt is not handled by the node or one of it's sub nodes, the interrupt is dropped.
//
// To delay the handling of an interrupt, use a Timer and call SendInterrupt from the HandleTimer method.
func (n *AbstractNode) SendInterrupt(ctx context.Context, interrupt Interrupt, to Address) {
	select {
	case <-ctx.Done():
		return
	default:
		n.sim.LogSendInterrupt(n.address, to, interrupt)
		it := InterruptTriplet{interrupt, n.address, to}
		if handled := n.sim.HandleInterrupt(ctx, it); !handled {
			n.sim.DropInterrupt(ctx, it)
		}
	}
}

// HandleInterrupt handles an interrupt received by the node.
//
// If the interrupt is a StopInterrupt, the node is stopped and cannot be started again.
//
// If the interrupt is a SleepInterrupt, the node is put to sleep for the specified duration and resumed after.
//
// If the interrupt is a StartInterrupt, the node is resumed, usually after sleeping for a specified duration.
//
// If an unknown interrupt is received, the function returns false, otherwise true.
func (n *AbstractNode) HandleInterrupt(ctx context.Context, interrupt Interrupt, from Address) bool {
	switch interrupt.Type {
	case StopInterrupt:
		n.state = Stopped
		return true
	case SleepInterrupt:
		data := interrupt.Data.(SleepInterruptData)
		n.state = Sleeping
		go func() {
			<-time.After(data.Duration)
			startInterrupt := NewInterrupt(StartInterrupt, nil)
			n.SendInterrupt(ctx, startInterrupt, n.address)
		}()
		return true
	case StartInterrupt:
		n.state = Running
		return true
	default:
		return false
	}
}
