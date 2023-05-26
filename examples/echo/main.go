package main

import (
	"time"

	ds "github.com/samuel-adekunle/disse"
)

func main() {
	// opts contains the configuration options for the simulation.
	// Either all fields are set or pass nil to use the default values.
	opts := &ds.SimulationOptions{
		Duration:     10 * time.Second,
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

	// TODO: 3. Create an echo node and add it to the simulation.
	echoAddress := ds.Address("echo")

	// Create a hello node and add it to the simulation.
	helloAddress := ds.Address("hello")
	helloNode := &HelloNode{
		AbstractNode: ds.NewAbstractNode(sim, helloAddress),
		receiver:     echoAddress,
	}
	sim.AddNode(helloAddress, helloNode)

	// TODO: 4. Run the simulation.
}
