package main

import (
	ds "disse/lib"
	"time"
)

func main() {
	nodes := make(map[ds.Address]ds.Node)
	mq := make(chan ds.MessageTriplet, 10)
	tq := make(chan ds.TimerTriplet, 10)

	ping, pong := ds.Message("Ping"), ds.Message("Pong")
	pingCounter := 3
	pingInterval := 100 * time.Millisecond

	sa := ds.Address("PingServer")
	nodes[sa] = &PingServer{ds.BaseNode{Address: sa, MessageQueue: mq, TimerQueue: tq}, ping, pong, 0}

	ca := ds.Address("PingClient")
	nodes[ca] = &PingClient{ds.BaseNode{Address: ca, MessageQueue: mq, TimerQueue: tq}, ping, pong, sa, pingInterval, pingCounter, 0}

	minLatency := 10 * time.Millisecond
	maxLatency := 50 * time.Millisecond
	simDuration := 3 * time.Second
	sim := ds.Simulation{
		Nodes:        nodes,
		MessageQueue: mq,
		TimerQueue:   tq,
		MinLatency:   minLatency,
		MaxLatency:   maxLatency,
		Duration:     simDuration,
	}
	sim.Run()
}
