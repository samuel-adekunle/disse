package disse

import (
	"context"
	"sync"
	"time"
)

// Node is the interface that must be implemented by all nodes in the distributed system.
type Node interface {
	// To be implemented by the user
	Init(context.Context)
	HandleMessage(context.Context, Message, Address) (handled bool)
	HandleTimer(context.Context, Timer, time.Duration) (handled bool)

	// To be implemented by the simulation engine i.e. LocalNode and LocalSimulation
	GetAddress() Address
	GetState() NodeState
	AddSubNode(Node)
	RemoveSubNode(Address)
	SendMessage(context.Context, Message, Address)
	BroadcastMessage(context.Context, Message, []Address)
	SetTimer(context.Context, Timer, time.Duration)
	SendInterrupt(context.Context, Interrupt, Address)

	// Private implementation specific functions
	initSubNodes(context.Context)
	findMessageHandler(context.Context, MessageTriplet) (handled bool)
	findTimerHandler(context.Context, TimerTriplet) (handled bool)
	findInterruptHandler(context.Context, InterruptTriplet) (handled bool)
	handleInterrupt(context.Context, Interrupt, Address) (handled bool)
}

// NodeState is a string that represents the state of a node.
//
// It can be either Stopped, Running or Sleeping.
type NodeState string

const (
	// Running is the state of a node that is running.
	Running NodeState = "Running"
	// Sleeping is the state of a node that is sleeping.
	Sleeping NodeState = "Sleeping"
	// Stopped is the state of a node that is stopped.
	Stopped NodeState = "Stopped"
)

// LocalNode implements most of the functions needed to satisfy the INode interface.
//
// It is designed to be used with LocalSimulation and should be used as a base for custom nodes.
type LocalNode struct {
	address  Address
	sim      *LocalSimulation
	subNodes map[Address]Node
	state    NodeState
}

// NewLocalNode creates a new LocalNode.
func NewLocalNode(sim *LocalSimulation, address Address) *LocalNode {
	return &LocalNode{
		address:  address,
		sim:      sim,
		subNodes: make(map[Address]Node),
		state:    Running,
	}
}

// GetAddress returns the address of the node.
func (n *LocalNode) GetAddress() Address {
	return n.address
}

// GetState returns the state of the node.
func (n *LocalNode) GetState() NodeState {
	return n.state
}

// AddSubNode adds a sub node to the node.
func (n *LocalNode) AddSubNode(node Node) {
	address := node.GetAddress()
	n.subNodes[address] = node
}

// RemoveSubNode removes a sub node with the given address.
func (n *LocalNode) RemoveSubNode(address Address) {
	delete(n.subNodes, address)
}

// initSubNodes initializes all sub nodes of a node.
func (n *LocalNode) initSubNodes(ctx context.Context) {
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			n.sim.initNode(ctx, _node)
		}(node)
	}
	wg.Wait()
}

// findMessageHandler finds the correct sub node to handle a message.
//
// If the message is for the current sub node, the HandleMessage method is called.
//
// If the message is not for the current sub node, the findMessageHandler method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the message is dropped.
//
// If a match is found, and the message is handled successfully, the method returns true, otherwise it returns false.
func (n *LocalNode) findMessageHandler(ctx context.Context, mt MessageTriplet) (handled bool) {
	if node, ok := n.subNodes[mt.To]; ok {
		return n.sim._handleMessage(ctx, node, mt)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.findMessageHandler(ctx, mt)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
}

// findTimerHandler finds the correct sub node to handle a timer.
//
// If the timer is for the current sub node, the HandleTimer method is called.
//
// If the timer is not for the current sub node, the findTimerHandler method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the timer is dropped.
//
// If a match is found, and the timer is handled successfully, the method returns true, otherwise it returns false.
func (n *LocalNode) findTimerHandler(ctx context.Context, tt TimerTriplet) (handled bool) {
	if node, ok := n.subNodes[tt.To]; ok {
		return n.sim._handleTimer(ctx, node, tt)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.findTimerHandler(ctx, tt)
			wg.Done()
		}(node)
	}
	wg.Wait()
	return handled
}

// findInterruptHandler handles an interrupt for all sub nodes.
//
// If the interrupt is for the node, the HandleInterrupt method is called.
//
// If the interrupt is not for the node, the findInterruptHandler method is called recursively to check it's sub nodes for a match.
//
// If no match is found or the matching node is not running, the interrupt is dropped.
//
// If an unknown interrupt is received, the interrupt is dropped and the function returns false, otherwise true.
func (n *LocalNode) findInterruptHandler(ctx context.Context, it InterruptTriplet) (handled bool) {
	if node, ok := n.subNodes[it.To]; ok {
		return n.sim._handleInterrupt(ctx, node, it)
	}
	var wg sync.WaitGroup
	for _, node := range n.subNodes {
		wg.Add(1)
		go func(_node Node) {
			handled = handled || _node.findInterruptHandler(ctx, it)
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
func (n *LocalNode) SendMessage(ctx context.Context, message Message, to Address) {
	select {
	case <-ctx.Done():
		return
	default:
		n.sim.LogSendMessage(n.address, to, message)
		mt := MessageTriplet{message, n.address, to}
		go func() {
			if to.Root() != n.address.Root() {
				time.Sleep(n.sim.randomLatency())
			}
			n.sim.messageQueue[mt.To] <- mt
		}()
	}
}

// BroadcastMessage sends a message to multiple nodes in the simulation.
//
// See SendMessage for more details on how the message is sent.
func (n *LocalNode) BroadcastMessage(ctx context.Context, message Message, to []Address) {
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
func (n *LocalNode) SetTimer(ctx context.Context, timer Timer, duration time.Duration) {
	select {
	case <-ctx.Done():
		return
	default:
		n.sim.LogSetTimer(n.address, timer, duration)
		go func() {
			time.Sleep(duration)
			n.sim.timerQueue[n.address] <- TimerTriplet{timer, n.address, duration}
		}()
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
func (n *LocalNode) SendInterrupt(ctx context.Context, interrupt Interrupt, to Address) {
	select {
	case <-ctx.Done():
		return
	default:
		n.sim.LogSendInterrupt(n.address, to, interrupt)
		it := InterruptTriplet{interrupt, n.address, to}
		go func() {
			n.sim.interruptQueue[it.To] <- it
		}()
	}
}

// handleInterrupt handles an interrupt received by the node.
//
// If the interrupt is a StopInterrupt, the node is stopped and cannot be started again.
//
// If the interrupt is a SleepInterrupt, the node is put to sleep for the specified duration and resumed after.
//
// If the interrupt is a StartInterrupt, the node is resumed, usually after sleeping for a specified duration.
//
// If an unknown interrupt is received, the function returns false, otherwise true.
func (n *LocalNode) handleInterrupt(ctx context.Context, interrupt Interrupt, from Address) bool {
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
