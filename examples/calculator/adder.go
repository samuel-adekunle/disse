package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

// AdderNode is a node that adds numbers.
type AdderNode struct {
	*ds.LocalNode
}

// Init is called when the node is initialized by the simulation.
func (n *AdderNode) Init(ctx context.Context) {}

// HandleMessage is called when the node receives a message.
func (n *AdderNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case CalculatorAdd:
		data := message.Data.(CalculatorOperationData)
		result := data.A + data.B
		resultMessage := ds.NewMessage(CalculatorResult, CalculatorResultData{
			A:         data.A,
			B:         data.B,
			Operation: CalculatorAdd,
			Result:    result,
		})
		n.SendMessage(ctx, resultMessage, from)
		return true
	case CalculatorSubtract:
		data := message.Data.(CalculatorOperationData)
		result := data.A - data.B
		resultMessage := ds.NewMessage(CalculatorResult, CalculatorResultData{
			A:         data.A,
			B:         data.B,
			Operation: CalculatorSubtract,
			Result:    result,
		})
		n.SendMessage(ctx, resultMessage, from)
		return true
	default:
		return false
	}
}

// HandleTimer is called when a node receives a timer.
func (n *AdderNode) HandleTimer(ctx context.Context, timer ds.Timer, duration time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
