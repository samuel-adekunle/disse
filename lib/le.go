package lib

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
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
// It assumes a crash-stop process abstraction and isn't useful for crash-recovery
// or Byzantine process abstractions.
//
// It elects a new leader when the current leader crashes and uses a perfect failure detector
// to detect crashes.
type LeNode struct {
	*ds.AbstractNode
	// External perfect failure detector
	Pfd     *PfdNode
	nodes   []ds.Address
	leader  ds.Address
	crashed map[ds.Address]bool
}

// Init is called when the node is initialized by the simulation.
func (n *LeNode) Init(ctx context.Context) {
	n.leader = n.nodes[0]
	n.crashed = make(map[ds.Address]bool)
	for _, node := range n.nodes {
		n.crashed[node] = false
	}
}

// HandleMessage is called when the node receives a message.
//
// If the message is a crash message, the node is marked as crashed.
// If the current leader is crashed, a new leader is elected.
func (n *LeNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case PfdCrash:
		data := message.Data.(PfdCrashData)
		n.crashed[data.Node] = true
		if n.leader == data.Node {
			for _, node := range n.nodes {
				if !n.crashed[node] {
					n.leader = node
					break
				}
			}
			// HACK(samuel-adekunle): Assume that a new leader is always elected.
			leaderMessage := ds.NewMessage(LeLeader, LeLeaderData{Node: n.leader})
			for _, node := range n.nodes {
				if !n.crashed[node] {
					n.SendMessage(ctx, leaderMessage, node)
				}
			}
		}
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
