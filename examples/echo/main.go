package main

import (
	"fmt"

	ds "github.com/samuel-adekunle/disse"
)

func main() {
	sim := ds.NewLocalSimulation(nil)

	echoAddress := ds.Address("echo")
	echoNode := &EchoNode{
		LocalNode: ds.NewLocalNode(sim, echoAddress),
	}
	sim.AddNode(echoNode)

	for i := 0; i < 3; i++ {
		helloAddress := ds.Address(fmt.Sprintf("hello%d", i))
		helloNode := &HelloNode{
			LocalNode: ds.NewLocalNode(sim, helloAddress),
			echoNode:  echoAddress,
		}
		sim.AddNode(helloNode)
	}

	sim.Run()
}
