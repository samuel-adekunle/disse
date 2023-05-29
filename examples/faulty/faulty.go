package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
	"github.com/samuel-adekunle/disse/lib"
)

const (
	// FaultyTimer is the type of timer used to crash a faulty node after 1 second.
	FaultyTimer ds.TimerType = "FaultyTimer"
)

// FaultyNode is a node that crashes after a given duration.
type FaultyNode struct {
	*ds.LocalNode
	lifetime time.Duration
}

// Init is called when the node is initialized by the simulation.
func (n *FaultyNode) Init(ctx context.Context) {
	faultyTimer := ds.NewTimer(FaultyTimer, nil)
	n.SetTimer(ctx, faultyTimer, n.lifetime)
}

// HandleMessage is called when the node receives a message.
func (n *FaultyNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case lib.PfdCrash:
		data := message.Data.(lib.PfdCrashData)
		fmt.Printf("%s received PfdCrash: %v\n", n.GetAddress(), data)
		return true
	case lib.LeLeader:
		data := message.Data.(lib.LeLeaderData)
		fmt.Printf("%s received LeLeader: %v\n", n.GetAddress(), data)
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
func (n *FaultyNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	case FaultyTimer:
		stopInterrupt := ds.NewInterrupt(ds.StopInterrupt, nil)
		n.SendInterrupt(ctx, stopInterrupt, n.GetAddress())
		return true
	default:
		return false
	}
}
