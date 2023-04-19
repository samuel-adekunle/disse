package main

import (
	"flag"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

var debugLogPath string
var umlLogPath string

// init sets up the command line flags for the simulation executable.
//
// The log file name is the file where the simulation logs will be written.
//
// The UML file name is the file where the UML diagram will be written.
//
// By default both logs are disabled (by setting their output to /dev/null) and the UML diagram is not generated.
func init() {
	defaultLogPath := "/dev/null"
	logFileNameUsage := "path to log file"
	defaultUmlPath := "/dev/null"
	umlFileNameUsage := "path to UML diagram file"

	flag.StringVar(&debugLogPath, "logfile", defaultLogPath, "path to log file")
	flag.StringVar(&debugLogPath, "l", defaultLogPath, logFileNameUsage+" (shorthand)")
	flag.StringVar(&umlLogPath, "umlfile", defaultUmlPath, "path to UML diagram file")
	flag.StringVar(&umlLogPath, "u", defaultUmlPath, umlFileNameUsage+" (shorthand)")
}

var sim *ds.Simulation
var beb *BebNode
var helloNodes []ds.Node
var faultyHelloNode *FaultyHelloNode

// initBebSimulation initializes the simulation.
func initBebSimulation() {
	bebAddress := ds.Address("Beb")
	helloNodeAddresses := []ds.Address{
		ds.Address("HelloNode1"),
		ds.Address("HelloNode2"),
		ds.Address("HelloNode3"),
	}

	sim = ds.NewSimulation(&ds.SimulationOptions{
		MinLatency:        ds.DefaultMinLatency,
		MaxLatency:        ds.DefaultMaxLatency,
		MessageBufferSize: ds.DefaultMessageBufferSize,
		TimerBufferSize:   ds.DefaultTimerBufferSize,
		DebugLogPath:      debugLogPath,
		UmlLogPath:        umlLogPath,
		Duration:          5 * time.Second,
	})
	beb = &BebNode{
		AbstractNode: ds.NewAbstractNode(sim, bebAddress),
		nodes:        append(helloNodeAddresses, bebAddress),
	}
	helloNodes = make([]ds.Node, 0)
	for i, helloNodeAddress := range helloNodeAddresses {
		helloNodes = append(helloNodes, &HelloNode{
			AbstractNode: ds.NewAbstractNode(sim, helloNodeAddress),
			beb:          bebAddress,
			sendAfter:    time.Second * time.Duration(i),
		})
	}
	sim.AddNode(bebAddress, beb)
	sim.AddNodes(helloNodeAddresses, helloNodes)
}

func main() {
	flag.Parse()
	initBebSimulation()
	sim.Run()
}
