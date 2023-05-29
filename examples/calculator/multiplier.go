package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

// MultiplierNode is a node that multiplies numbers.
type MultiplierNode struct {
	*ds.LocalNode
}

// Init is called when the node is initialized by the simulation.
func (n *MultiplierNode) Init(ctx context.Context) {}

// HandleMessage is called when the node receives a message.
func (n *MultiplierNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case CalculatorMultiply:
		data := message.Data.(CalculatorOperationData)
		result := data.A * data.B
		resultMessage := ds.NewMessage(CalculatorResult, CalculatorResultData{
			A:         data.A,
			B:         data.B,
			Operation: CalculatorMultiply,
			Result:    result,
		})
		n.SendMessage(ctx, resultMessage, from)
		return true
	case CalculatorDivide:
		data := message.Data.(CalculatorOperationData)
		result := data.A / data.B
		resultMessage := ds.NewMessage(CalculatorResult, CalculatorResultData{
			A:         data.A,
			B:         data.B,
			Operation: CalculatorDivide,
			Result:    result,
		})
		n.SendMessage(ctx, resultMessage, from)
		return true
	default:
		return false
	}
}

// HandleTimer is called when a node receives a timer.
func (n *MultiplierNode) HandleTimer(ctx context.Context, timer ds.Timer, duration time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
