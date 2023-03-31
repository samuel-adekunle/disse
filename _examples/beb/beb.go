package main

import (
	"time"

	ds "github.com/samuel-adekunle/disse"
)

var sim *ds.Simulation
var bebServer *BebServer
var bebClients []ds.Node

func initBebSimulation() {
	bebServerAddress := ds.Address("BebServer")
	bebClientAddresses := []ds.Address{
		ds.Address("BebClient1"),
		ds.Address("BebClient2"),
		ds.Address("BebClient3"),
	}

	sim = ds.NewSimulation()
	bebServer = &BebServer{
		AbstractNode: ds.NewAbstractNode(sim, bebServerAddress),
		nodes:        bebClientAddresses,
	}
	bebClients = make([]ds.Node, 0)
	for i, clientAddress := range bebClientAddresses {
		bebClients = append(bebClients, &BebClient{
			AbstractNode: ds.NewAbstractNode(sim, clientAddress),
			bebServer:    bebServerAddress,
			messageDelay: time.Second * time.Duration(i),
		})
	}
	sim.AddNode(bebServerAddress, bebServer)
	sim.AddNodes(bebClientAddresses, bebClients)
	sim.Duration = 5 * time.Second
}

func init() {
	initBebSimulation()
}

func main() {
	sim.Run()
}
