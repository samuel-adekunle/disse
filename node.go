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
	// Implemented by BaseNode:

	GetAddress() Address
	GetState() NodeState
	AddSubNode(Address, Node)
	SubNodesInit(context.Context)
	SubNodesHandleMessage(context.Context, MessageTriplet) (handled bool)
	SubNodesHandleTimer(context.Context, TimerTriplet) (handled bool)
	SubNodesHandleInterrupt(context.Context, InterruptTriplet) (handled bool)
	HandleInterrupt(context.Context, Interrupt) (handled bool)
	SendMessage(context.Context, Message, Address)
	BroadcastMessage(context.Context, Message, []Address)
	SetTimer(context.Context, Timer, time.Duration)
	SendInterrupt(context.Context, Interrupt, Address)

	// To be implemented by concrete node:

	Init(context.Context)
	HandleMessage(context.Context, Message, Address) (handled bool)
	HandleTimer(context.Context, Timer, time.Duration) (handled bool)
}

// NodeState is a string that represents the state of a node.
//
// It can be either Stopped, Running or Sleeping.
type NodeState string

const (
	Stopped  NodeState = "Stopped"
	Running  NodeState = "Running"
	Sleeping NodeState = "Sleeping"
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

// AddSubNode adds a sub node to the node.
func (n *AbstractNode) AddSubNode(address Address, node Node) {
	n.subNodes[address] = node
}

// SubNodesInit initializes all sub nodes.
func (n *AbstractNode) SubNodesInit(ctx context.Context) {
	var wg sync.WaitGroup
	for address, node := range n.subNodes {
		wg.Add(1)
		go func(_address Address, _node Node) {
			n.sim.debugLog.LogNodeState(_node)
			_node.Init(ctx)
			_node.SubNodesInit(ctx)
			wg.Done()
		}(address, node)
	}
	wg.Wait()
}

// SubNodesHandleMessage handles a message for all sub nodes.
//
// If the message is for the node, the HandleMessage method is called.
//
// If the message is not for the node, the SubNodesHandleMessage method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the message is dropped.
//
// If a match is found, and the message is handled successfully, the method returns true, otherwise it returns false.
func (n *AbstractNode) SubNodesHandleMessage(ctx context.Context, mt MessageTriplet) (handled bool) {
	if node, ok := n.subNodes[mt.To]; ok {
		if node.GetState() != Running {
			return false
		}
		n.sim.debugLog.LogHandleMessage(mt.From, mt.To, mt.Message)
		return node.HandleMessage(ctx, mt.Message, mt.From)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.SubNodesHandleMessage(ctx, mt)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
}

// SubNodesHandleTimer handles a timer for all sub nodes.
//
// If the timer is for the node, the HandleTimer method is called.
//
// If the timer is not for the node, the SubNodesHandleTimer method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the timer is dropped.
//
// If a match is found, the timer is successfully handled, the method returns true, otherwise it returns false.
func (n *AbstractNode) SubNodesHandleTimer(ctx context.Context, tt TimerTriplet) (handled bool) {
	if node, ok := n.subNodes[tt.To]; ok {
		if node.GetState() != Running {
			return false
		}
		n.sim.debugLog.LogHandleTimer(tt.To, tt.Timer, tt.Duration)
		return node.HandleTimer(ctx, tt.Timer, tt.Duration)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.SubNodesHandleTimer(ctx, tt)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
}

// SubNodesHandleInterrupt handles an interrupt for all sub nodes.
//
// If the interrupt is for the node, the HandleInterrupt method is called.
//
// If the interrupt is not for the node, the SubNodesHandleInterrupt method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the interrupt is dropped.
//
// If an unknown interrupt is received, the interrupt is dropped and the function returns false, otherwise true.
func (n *AbstractNode) SubNodesHandleInterrupt(ctx context.Context, ip InterruptTriplet) (handled bool) {
	if node, ok := n.subNodes[ip.To]; ok {
		if node.GetState() == Stopped {
			return false
		}
		n.sim.debugLog.LogHandleInterrupt(ip.From, ip.To, ip.Interrupt)
		return node.HandleInterrupt(ctx, ip.Interrupt)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.SubNodesHandleInterrupt(ctx, ip)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
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
func (n *AbstractNode) HandleInterrupt(ctx context.Context, interrupt Interrupt) bool {
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
		n.sim.umlLog.LogSendMessage(n.address, to, message)
		n.sim.debugLog.LogSendMessage(n.address, to, message)
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
		n.sim.umlLog.LogSetTimer(n.address, timer, duration)
		n.sim.debugLog.LogSetTimer(n.address, timer, duration)
		n.sim.timerQueue <- TimerTriplet{timer, n.address, duration}
	}
}

// SendInterrupt sends an interrupt to another node in the simulation.
//
// Interrupts are handled immediately and do not go through the message queue.
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
		n.sim.umlLog.LogSendInterrupt(n.address, to, interrupt)
		n.sim.debugLog.LogSendInterrupt(n.address, to, interrupt)
		ip := InterruptTriplet{interrupt, n.address, to}
		if handled := n.sim.HandleInterrupt(ctx, ip); !handled {
			n.sim.DropInterrupt(ctx, ip)
		}
	}
}
