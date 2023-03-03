package main

import (
	ds "disse/lib"
	"time"
)

type PingServer struct {
	ds.BaseNode
	PingMessage ds.Message
	PongMessage ds.Message
	PingCounter int
}

func (n *PingServer) Init() {
	n.LogInit()
}

func (n *PingServer) HandleMessage(message ds.Message, from ds.Address) {
	n.BaseNode.LogHandleMessage(message, from)
	if message == n.PingMessage {
		n.PingCounter++
		n.SendMessage(n.PongMessage, from)
	}
}

func (n *PingServer) HandleTimer(timer ds.Timer, length time.Duration) {
	n.BaseNode.LogHandleTimer(timer, length)
}
