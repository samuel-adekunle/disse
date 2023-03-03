package main

import (
	ds "disse/lib"
	"time"
)

const PingTimer ds.Timer = "PingTimer"

type PingClient struct {
	ds.BaseNode
	PingMessage   ds.Message
	PongMessage   ds.Message
	ServerAddress ds.Address
	PingInterval  time.Duration
	PingCounter   int
	PongCounter   int
}

func (n *PingClient) Init() {
	n.LogInit()
	n.SetTimer(PingTimer, n.PingInterval)
}

func (n *PingClient) HandleMessage(message ds.Message, from ds.Address) {
	n.BaseNode.LogHandleMessage(message, from)
	if message == n.PongMessage {
		n.PongCounter++
	}
}

func (n *PingClient) HandleTimer(timer ds.Timer, length time.Duration) {
	n.BaseNode.LogHandleTimer(timer, length)
	if timer == PingTimer && n.PingCounter > 0 {
		n.SendMessage(n.PingMessage, n.ServerAddress)
		n.SetTimer(PingTimer, n.PingInterval)
		n.PingCounter--
	}
}
