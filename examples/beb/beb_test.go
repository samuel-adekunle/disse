package beb

import (
	ds "disse/lib"
	"testing"
	"time"
)

var sim *ds.Simulation
var bebServer *BebServer
var bebClients []ds.Node

func newSim() {
	bebServerAddress := ds.Address("BebServer")
	bebClientAddresses := []ds.Address{
		ds.Address("BebClient1"),
		ds.Address("BebClient2"),
		ds.Address("BebClient3"),
	}

	sim = ds.NewSimulation()
	bebServer = &BebServer{
		BaseNode: ds.NewBaseNode(sim, bebServerAddress),
		nodes:    bebClientAddresses,
	}
	bebClients = make([]ds.Node, 0)
	for i, clientAddress := range bebClientAddresses {
		bebClients = append(bebClients, &BebClient{
			BaseNode:     ds.NewBaseNode(sim, clientAddress),
			bebServer:    bebServerAddress,
			messageDelay: time.Second * time.Duration(i),
		})
	}
	sim.AddNode(bebServerAddress, bebServer)
	sim.AddNodes(bebClientAddresses, bebClients)
	sim.Duration = 5 * time.Second
}

func TestBeb(t *testing.T) {
	newSim()
	sim.Run()
}
