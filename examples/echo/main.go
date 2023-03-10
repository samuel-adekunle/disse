package main

import (
	ds "disse/lib"
	"time"
)

func main() {
	sim := ds.NewSimulation()
	echoServerAddress, echoClientAddress := ds.Address("EchoServer"), ds.Address("EchoClient")
	pingServerAddress, pingClientAddress := echoServerAddress.SubAddress("PingServer"), echoClientAddress.SubAddress("PingClient")
	pingMessage, pongMessage := ds.Message("Ping"), ds.Message("Pong")

	pingServer := &PingServer{
		BaseNode: ds.BaseNode{
			Address:      pingServerAddress,
			MessageQueue: sim.MessageQueue,
			TimerQueue:   sim.TimerQueue,
		},
		PingMessage: pingMessage,
		PongMessage: pongMessage,
		PingCounter: 0,
	}
	pingClient := &PingClient{
		BaseNode: ds.BaseNode{
			Address:      pingClientAddress,
			MessageQueue: sim.MessageQueue,
			TimerQueue:   sim.TimerQueue,
		},
		PingMessage:   pingMessage,
		PongMessage:   pongMessage,
		ServerAddress: pingServerAddress,
		PingInterval:  1 * time.Second,
		PongCounter:   0,
	}
	echoSever := &EchoServer{
		BaseNode: ds.BaseNode{
			Address:      echoServerAddress,
			MessageQueue: sim.MessageQueue,
			TimerQueue:   sim.TimerQueue,
		},
		PingServer: pingServer,
	}
	echoClient := &EchoClient{
		BaseNode: ds.BaseNode{
			Address:      echoClientAddress,
			MessageQueue: sim.MessageQueue,
			TimerQueue:   sim.TimerQueue,
		},
		PingClient: pingClient,
	}

	sim.AddNode(echoServerAddress, echoSever)
	sim.AddNode(echoClientAddress, echoClient)
	sim.Duration = 3 * time.Second
	sim.Run()
}
