package main

import (
	ds "disse/lib"
	"time"
)

func main() {
	echoServerAddress, echoClientAddress := ds.Address("EchoServer"), ds.Address("EchoClient")
	pingServerAddress, pingClientAddress := echoServerAddress.SubAddress("PingServer"), echoClientAddress.SubAddress("PingClient")
	pingMessage, pongMessage, echoMessage := ds.Message("Ping"), ds.Message("Pong"), ds.Message("Echo")

	sim := ds.NewSimulation()
	pingServer := &PingServer{
		BaseNode:    ds.NewBaseNode(sim, pingServerAddress),
		PingMessage: pingMessage,
		PongMessage: pongMessage,
		PingCounter: 0,
	}
	pingClient := &PingClient{
		BaseNode:      ds.NewBaseNode(sim, pingClientAddress),
		PingMessage:   pingMessage,
		PongMessage:   pongMessage,
		ServerAddress: pingServerAddress,
		PingInterval:  1 * time.Second,
		PongCounter:   0,
	}
	echoSever := &EchoServer{
		BaseNode:    ds.NewBaseNode(sim, echoServerAddress),
		EchoMessage: echoMessage,
		EchoCounter: 0,
	}
	echoClient := &EchoClient{
		BaseNode:          ds.NewBaseNode(sim, echoClientAddress),
		EchoInterval:      2 * time.Second,
		EchoServerAddress: echoServerAddress,
		EchoMessage:       echoMessage,
		EchoCounter:       0,
	}

	echoClient.AddSubNode(pingClientAddress, pingClient)
	echoSever.AddSubNode(pingServerAddress, pingServer)
	sim.AddNode(echoServerAddress, echoSever)
	sim.AddNode(echoClientAddress, echoClient)
	sim.Duration = 10 * time.Second
	sim.Run()
}
