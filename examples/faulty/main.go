package main

import (
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
	"github.com/samuel-adekunle/disse/lib"
)

func main() {
	sim := ds.NewLocalSimulation(nil)

	nodes := []ds.Address{}
	for i := 0; i < 5; i++ {
		faultyAddress := ds.Address(fmt.Sprintf("faulty%d", i))
		faultyNode := &FaultyNode{
			LocalNode: ds.NewLocalNode(sim, faultyAddress),
			lifetime:  time.Duration(i+1) * 2 * time.Second,
		}
		sim.AddNode(faultyNode)
		nodes = append(nodes, faultyAddress)
	}

	pfdAddress := ds.Address("pfd")
	nodes = append(nodes, pfdAddress)
	leAddress := ds.Address("le")
	leNode := &lib.LeNode{
		LocalNode: ds.NewLocalNode(sim, leAddress),
		Nodes:     nodes,
	}
	sim.AddNode(leNode)

	nodes = append(nodes, leAddress)
	pfdNode := &lib.PfdNode{
		LocalNode:       ds.NewLocalNode(sim, pfdAddress),
		Nodes:           nodes,
		TimeoutDuration: 500 * time.Millisecond,
	}
	sim.AddNode(pfdNode)

	sim.Run()
}
