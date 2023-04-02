package main

import (
	"time"

	ds "github.com/samuel-adekunle/disse"
)

var sim *ds.Simulation
var beb *BebNode
var helloNodes []ds.Node

// initBebSimulation initializes the simulation.
func initBebSimulation() {
	bebAddress := ds.Address("Beb")
	helloNodeAddresses := []ds.Address{
		ds.Address("HelloNode1"),
		ds.Address("HelloNode2"),
		ds.Address("HelloNode3"),
	}

	sim = ds.NewSimulation()
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
	sim.Duration = 5 * time.Second
}

func init() {
	initBebSimulation()
}

func main() {
	sim.Run()
}
