package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// FaultTimer is the type of timer used to stop a node after a given delay.
	FaultTimer ds.TimerType = "FaultTimer"
)

// FaultyNode is a node that stops itself after a delay.
type FaultyNode struct {
	*ds.AbstractNode
	interruptDelay time.Duration
}

// Init is called when the node is initialized by the simulation.
func (n *FaultyNode) Init(ctx context.Context) {
	faultTimer := ds.NewTimer(FaultTimer, nil)
	n.SetTimer(ctx, faultTimer, n.interruptDelay)
}

// HandleMessage is called when the node receives a message.
func (n *FaultyNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case PfdCrash:
		data := message.Data.(PfdCrashData)
		fmt.Printf("%s received PfdCrash: %v\n", n.GetAddress(), data)
		return true
	case PfdHeartbeatRequest:
		heartbeatReply := ds.NewMessage(PfdHeartbeatReply, nil)
		n.SendMessage(ctx, heartbeatReply, from)
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *FaultyNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	case FaultTimer:
		faultInterrupt := ds.NewInterrupt(ds.StopInterrupt, nil)
		n.SendInterrupt(ctx, faultInterrupt, n.GetAddress())
		return true
	default:
		return false
	}
}

func main() {
	// opts contains the configuration options for the simulation.
	// Either all fields are set or pass nil to use the default values.
	opts := &ds.SimulationOptions{
		Duration:          8 * time.Second,
		MinLatency:        10 * time.Millisecond,
		MaxLatency:        100 * time.Millisecond,
		MessageBufferSize: 100,
		TimerBufferSize:   100,
		DebugLogPath:      "debug.log",
		UmlLogPath:        "uml.log",
		JavaPath:          "/usr/bin/java",
		PlantumlPath:      "/usr/share/plantuml/plantuml.jar",
	}

	// sim is the simulation created with the given options.
	sim := ds.NewSimulation(opts)
	// nodes is a list of all nodes in the simulation.
	nodes := []ds.Address{}

	// Create 3 faulty nodes that stop after 2, 4, and 6 seconds respectively
	// and add them to the simulation.
	for i := 1; i <= 3; i++ {
		faultyAddress := ds.Address(fmt.Sprintf("faulty%d", i))
		nodes = append(nodes, faultyAddress)
		faultyNode := &FaultyNode{
			AbstractNode:   ds.NewAbstractNode(sim, faultyAddress),
			interruptDelay: time.Duration(i*2) * time.Second,
		}
		sim.AddNode(faultyAddress, faultyNode)
	}

	// Create a pfd node and add it to the simulation.
	pfdAddress := ds.Address("pfd")
	nodes = append(nodes, pfdAddress)
	pfdNode := &PfdNode{
		AbstractNode:    ds.NewAbstractNode(sim, pfdAddress),
		nodes:           nodes,
		timeoutDuration: 10 * opts.MaxLatency,
	}
	sim.AddNode(pfdAddress, pfdNode)

	// Run the simulation.
	sim.Run()
}
