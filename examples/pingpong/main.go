package main

import (
	ds "disse/lib"
	"time"
)

func main() {
	sim := ds.NewSimulation()
	pingMessage, pongMessage := ds.Message("Ping"), ds.Message("Pong")
	serverAddress, clientAddress := ds.Address("PingServer"), ds.Address("PingClient")
	pingServer := &PingServer{
		BaseNode:    ds.NewBaseNode(sim, serverAddress),
		PingMessage: pingMessage,
		PongMessage: pongMessage,
		PingCounter: 0,
	}
	pingClient := &PingClient{
		BaseNode:      ds.NewBaseNode(sim, clientAddress),
		PingMessage:   pingMessage,
		PongMessage:   pongMessage,
		ServerAddress: serverAddress,
		PingInterval:  1 * time.Second,
		PongCounter:   0,
	}
	sim.AddNode(serverAddress, pingServer)
	sim.AddNode(clientAddress, pingClient)
	sim.Duration = 10 * time.Second
	sim.Run()
}
