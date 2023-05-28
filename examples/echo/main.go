package main

import (
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

func main() {
	// opts contains the configuration options for the simulation.
	// Either all fields are set or pass nil to use the default values.
	opts := &ds.SimulationOptions{
		Duration:     5 * time.Second,
		MinLatency:   10 * time.Millisecond,
		MaxLatency:   100 * time.Millisecond,
		BufferSize:   10,
		DebugLogPath: "debug.log",
		UmlLogPath:   "uml.log",
		JavaPath:     "/usr/bin/java",
		PlantumlPath: "/usr/share/plantuml/plantuml.jar",
	}

	// sim is the simulation created with the given options.
	sim := ds.NewSimulation(opts)

	// Create an echo node and add it to the simulation.
	echoAddress := ds.Address("echo")
	echoNode := &EchoNode{
		AbstractNode: ds.NewAbstractNode(sim, echoAddress),
	}
	sim.AddNode(echoAddress, echoNode)

	// Create a 3 hello nodes and add them to the simulation.
	for i := 0; i < 3; i++ {
		helloAddress := ds.Address(fmt.Sprintf("hello%d", i))
		helloNode := &HelloNode{
			AbstractNode: ds.NewAbstractNode(sim, helloAddress),
			echoNode:     echoAddress,
		}
		sim.AddNode(helloAddress, helloNode)
	}

	// Run the simulation.
	sim.Run()
}
