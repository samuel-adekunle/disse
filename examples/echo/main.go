package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// HelloTimer is the type of timer used to send a hello to an echo node after 1 second.
	HelloTimer ds.TimerType = "HelloTimer"
	// EchoSend is the type of message used to send a hello.
	Hello ds.MessageType = "Hello"
)

// HelloData is the data of a hello message.
type HelloData string

// HelloNode is a node that sends a hello to an echo node after 1 second.
type HelloNode struct {
	*ds.AbstractNode
	echoNode ds.Address
}

// Init is called when the node is initialized by the simulation.
func (n *HelloNode) Init(ctx context.Context) {
	timer := ds.NewTimer(HelloTimer, nil)
	n.SetTimer(ctx, timer, 1*time.Second)
}

// HandleMessage is called when the node receives a message.
func (n *HelloNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case Hello:
		data := message.Data.(HelloData)
		fmt.Printf("%s received Hello: %s\n", n.GetAddress(), data)
		return true
	case EchoDeliver:
		data := message.Data.(EchoDeliverData)
		fmt.Printf("%s received EchoDeliver: %s\n", n.GetAddress(), data)
		return true
	default:
		return false
	}
}

// HandleTimer is called when the node receives a timer.
func (n *HelloNode) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) bool {
	switch timer.Type {
	case HelloTimer:
		echoSendMessage := ds.NewMessage(EchoSend, EchoSendData{
			Message: ds.NewMessage(Hello, HelloData("Hello")),
		})
		n.SendMessage(ctx, echoSendMessage, n.echoNode)
		return true
	default:
		return false
	}
}

func main() {
	// opts contains the configuration options for the simulation.
	// Either all fields are set or pass nil to use the default values.
	opts := &ds.SimulationOptions{
		Duration:     3 * time.Second,
		MinLatency:   10 * time.Millisecond,
		MaxLatency:   100 * time.Millisecond,
		BufferSize:   20,
		DebugLogPath: "debug.log",
		UmlLogPath:   "uml.log",
		JavaPath:     "/usr/bin/java",
		PlantumlPath: "/usr/share/plantuml/plantuml.jar",
	}

	// sim is the simulation created with the given options.
	sim := ds.NewSimulation(opts)

	// TODO: Create an echo node and add it to the simulation.
	echoAddress := ds.Address("echo")

	// Create a hello node and add it to the simulation.
	helloAddress := ds.Address("hello")
	helloNode := &HelloNode{
		AbstractNode: ds.NewAbstractNode(sim, helloAddress),
		echoNode:     echoAddress,
	}
	sim.AddNode(helloAddress, helloNode)

	// Run the simulation.
	sim.Run()
}
