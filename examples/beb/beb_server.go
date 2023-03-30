package beb

import (
	"context"
	ds "disse/lib"
	"time"
)

type BebServer struct {
	*ds.BaseNode
	nodes []ds.Address
}

func (n *BebServer) Init(ctx context.Context) {}

func (n *BebServer) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) {

}

func (n *BebServer) HandleTimer(ctx context.Context, timer ds.Timer, length time.Duration) {}
