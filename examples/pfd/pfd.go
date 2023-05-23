package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// PfdTimeout is the type of timer used to detect process failures.
	PfdTimeout ds.TimerType = "PfdTimeout"
	// PfdCrash is the type of message used to indicate that a node has crashed.
	PfdCrash ds.MessageType = "PfdCrash"
	// PfdHeartbeatRequest is the type of message used to request a heartbeat.
	PfdHeartbeatRequest ds.MessageType = "PfdHeartbeatRequest"
	// PfdHeartbeatReply is the type of message used to reply to a heartbeat request.
	PfdHeartbeatReply ds.MessageType = "PfdHeartbeatReply"
)

// PfdCrashData is the data of a crash message.
type PfdCrashData struct {
	Node ds.Address
}

// PfdNode is a node that implements a perfect failure detector which assumes a crash-stop
// process abstraction and a synchronous system with a known upper bound on message delay.
//
// This implementation uses the "Exclude on Timeout" algorithm.
type PfdNode struct {
	*ds.AbstractNode
	Nodes           []ds.Address
	alive           map[ds.Address]bool
	crashed         map[ds.Address]bool
	timeoutDuration time.Duration
}

// Init is called when the node is initialized by the simulation.
func (n *PfdNode) Init(ctx context.Context) {
	// TODO: initialize alive and crashed maps
	// TODO: set timeout timer
}

// HandleMessage is called when the node receives a message.
func (n *PfdNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case PfdHeartbeatRequest:
		heartbeatReplyMessage := ds.NewMessage(PfdHeartbeatReply, nil)
		n.SendMessage(ctx, heartbeatReplyMessage, from)
		return true
	case PfdHeartbeatReply:
		n.alive[from] = true
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *PfdNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	case PfdTimeout:
		aliveNodes := []ds.Address{}
		for _, node := range n.Nodes {
			if n.alive[node] {
				aliveNodes = append(aliveNodes, node)
			}
		}
		// TODO: broadcast crash message for crashed nodes
		heartbeatRequestMessage := ds.NewMessage(PfdHeartbeatRequest, nil)
		n.BroadcastMessage(ctx, heartbeatRequestMessage, aliveNodes)
		for _, node := range n.Nodes {
			n.alive[node] = false
		}
		timeoutTimer := ds.NewTimer(PfdTimeout, nil)
		n.SetTimer(ctx, timeoutTimer, n.timeoutDuration)
		return true
	default:
		return false
	}
}
