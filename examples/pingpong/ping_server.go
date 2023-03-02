package main

import (
	ds "disse/lib"
)

type PingServer struct {
	ds.BaseNode
	PingMessage ds.Message
	PongMessage ds.Message
}

func (n *PingServer) Init() {}

func (n *PingServer) HandleMessage(message ds.Message, from ds.Address) {
	if message == n.PingMessage {
		n.SendMessage(n.PongMessage, from)
	}
}

func (n *PingServer) HandleTimer(timer ds.Timer) {}
