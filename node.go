package disse

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

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

// Node is the interface that must be implemented by all nodes in the distributed system.
type Node interface {
	Init(context.Context)
	GetAddress() Address
	GetState() NodeState
	GetSubNodes() map[Address]Node
	AddSubNode(Node) error
	RemoveSubNode(Address)
	SendMessage(context.Context, Message, Address) error
	BroadcastMessage(context.Context, Message, []Address) error
	SetTimer(context.Context, Timer, time.Duration) error
	SendInterrupt(context.Context, Interrupt, Address) error
	HandleMessage(context.Context, Message, Address) (handled bool)
	HandleTimer(context.Context, Timer, time.Duration) (handled bool)
	HandleInterrupt(context.Context, Interrupt, Address) (handled bool)
}

// LocalNode implements most of the functions needed to satisfy the INode interface.
//
// The only functions that need to be implemented by a LocalNode are Init, HandleMessage and HandleTimer.
//
// It is designed to be used with a LocalSimulation.
type LocalNode struct {
	address  Address
	sim      *LocalSimulation
	subNodes map[Address]Node
	state    NodeState
}

// NewLocalNode creates a new LocalNode with the given address.
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

// GetSubNodes returns the sub nodes of the node.
func (n *LocalNode) GetSubNodes() map[Address]Node {
	return n.subNodes
}

// AddSubNode adds a sub node to the node.
//
// Parent nodes need to be added to the simulation before adding sub nodes.
func (n *LocalNode) AddSubNode(node Node) error {
	address := node.GetAddress()
	if _, ok := n.subNodes[address]; ok {
		return fmt.Errorf("node with address %s already exists", address)
	}
	if err := n.validateNode(address.GetRoot()); err != nil {
		return err
	}
	n.subNodes[address] = node
	return nil
}

// RemoveSubNode removes a sub node with the given address.
func (n *LocalNode) RemoveSubNode(address Address) {
	delete(n.subNodes, address)
}

// SendMessage sends a message to another node in the simulation.
//
// A random amount of latency will be added to the message if the sender and receiver are not the same node.
//
// If the destination node is not valid, an error is returned.
func (n *LocalNode) SendMessage(ctx context.Context, message Message, to Address) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		if err := n.validateNode(to); err != nil {
			return err
		}
		from := n.address.GetRoot()
		n.sim.LogSendMessage(from, to, message)
		go func() {
			if to != from {
				time.Sleep(n.randomLatency())
			}
			n.sim.messageQueue[to] <- MessageTriplet{message, from, to}
		}()
		return nil
	}
}

// BroadcastMessage sends a message to multiple nodes in the simulation.
//
// See SendMessage for more details on how the message is sent.
//
// If any of the destination nodes  are not valid, an error is returned.
func (n *LocalNode) BroadcastMessage(ctx context.Context, message Message, to []Address) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		for _, address := range to {
			if err := n.validateNode(address); err != nil {
				return err
			}
		}
		for _, address := range to {
			n.SendMessage(ctx, message, address)
		}
		return nil
	}
}

// SetTimer sets a timer for the node.
//
// The timer is added to the timer queue of the node after the given duration.
//
// If the destination node is not valid, an error is returned.
func (n *LocalNode) SetTimer(ctx context.Context, timer Timer, duration time.Duration) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		to := n.address.GetRoot()
		if err := n.validateNode(to); err != nil {
			return err
		}
		from := to
		n.sim.LogSetTimer(to, timer, duration)
		go func() {
			time.Sleep(duration)
			n.sim.timerQueue[to] <- TimerTriplet{timer, from, duration}
		}()
		return nil
	}
}

// SendInterrupt sends an interrupt to another node in the simulation.
//
// Interrupts are always immediately added to the interrupt queue of the destination node.
//
// If the destination node does is not valid, an error is returned.
func (n *LocalNode) SendInterrupt(ctx context.Context, interrupt Interrupt, to Address) error {
	select {
	case <-ctx.Done():
		return nil
	default:
		if err := n.validateNode(to); err != nil {
			return err
		}
		from := n.address.GetRoot()
		n.sim.LogSendInterrupt(from, to, interrupt)
		go func() {
			n.sim.interruptQueue[to] <- InterruptTriplet{interrupt, from, to}
		}()
		return nil
	}
}

// HandleInterrupt handles an interrupt received by the node.
func (n *LocalNode) HandleInterrupt(ctx context.Context, interrupt Interrupt, from Address) bool {
	switch interrupt.Type {
	case StopInterrupt:
		n.state = Stopped
		return true
	case SleepInterrupt:
		data := interrupt.Data.(SleepInterruptData)
		n.state = Sleeping
		go func() {
			<-time.After(data.Duration)
			n.state = Running
		}()
		return true
	default:
		return false
	}
}

// randomLatency returns a random duration between the minimum and maximum latency.
func (n *LocalNode) randomLatency() time.Duration {
	s := n.sim
	return s.options.MinLatency + time.Duration(rand.Int63n(int64(s.options.MaxLatency-s.options.MinLatency)))
}

// validateNode checks if the node exists in the simulation.
func (n *LocalNode) validateNode(address Address) error {
	if _, ok := n.sim.nodes[address]; !ok {
		return fmt.Errorf("node with address %s does not exist", address)
	}
	return nil
}
