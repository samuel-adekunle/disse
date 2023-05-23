package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
	lib "github.com/samuel-adekunle/disse/lib"
)

const (
	// LeLeader is the type of message used to indicate that a node is the leader.
	LeLeader ds.MessageType = "LeLeader"
)

// LeLeaderData is the data of a leader message.
type LeLeaderData struct {
	Node ds.Address
}

// LeNode is a node that implements leader election.
//
// It elects a new leader when the current leader crashes
// using the "Monarchical Leader Election" algorithm where node
// index is used as the rank.
//
// It uses a perfect failure detector to detect crashes which assumes
// a crash-stop process abstraction and a synchronous system with a known
// upper bound on message delay.
type LeNode struct {
	*ds.AbstractNode
	Nodes   []ds.Address
	leader  ds.Address
	crashed map[ds.Address]bool
}

// Init is called when the node is initialized by the simulation.
func (n *LeNode) Init(ctx context.Context) {
	// TODO: initialize crashed map
	// TODO: set leader to node with lowest index
	// TODO: broadcast leader message
}

// HandleMessage is called when the node receives a message.
func (n *LeNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case // TODO: handle PfdCrash message from library:
		data := message.Data.(lib.PfdCrashData)
		n.crashed[data.Node] = true
		if n.leader != data.Node {
			return true
		}

		aliveNodes := []ds.Address{}
		for _, node := range n.Nodes {
			if !n.crashed[node] {
				aliveNodes = append(aliveNodes, node)
			}
		}
		if len(aliveNodes) == 0 {
			return true
		}
		// TODO: set leader to node with lowest index
		// TODO: broadcast leader message
		return true
	case lib.PfdHeartbeatRequest:
		heartbeatReply := ds.NewMessage(lib.PfdHeartbeatReply, nil)
		n.SendMessage(ctx, heartbeatReply, from)
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *LeNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
