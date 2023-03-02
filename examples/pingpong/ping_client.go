package main

import (
	ds "disse/lib"
	"time"
)

const PingTimer ds.Timer = "PingTimer"

type PingClient struct {
	ds.BaseNode
	PingCounter   int
	PingInterval  time.Duration
	ServerAddress ds.Address
	PingMessage   ds.Message
	PongMessage   ds.Message
}

func (n *PingClient) Init() {
	n.SetTimer(PingTimer, n.PingInterval)
}

func (n *PingClient) HandleMessage(message ds.Message, from ds.Address) {}

func (n *PingClient) HandleTimer(timer ds.Timer) {
	if timer == PingTimer && n.PingCounter > 0 {
		n.SendMessage(n.PingMessage, n.ServerAddress)
		n.SetTimer(PingTimer, n.PingInterval)
		n.PingCounter--
	}
}
