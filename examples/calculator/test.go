package main

import (
	"context"
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// SendAdd is the timer type for sending an add message.
	SendAdd ds.TimerType = "SendAdd"
	// SendSubtract is the timer type for sending a subtract message.
	SendSubtract ds.TimerType = "SendSubtract"
	// SendMultiply is the timer type for sending a multiply message.
	SendMultiply ds.TimerType = "SendMultiply"
	// SendDivide is the timer type for sending a divide message.
	SendDivide ds.TimerType = "SendDivide"
)

// TestNode is a node that tests the calculator node.
type TestNode struct {
	*ds.LocalNode
	A          int
	B          int
	calculator ds.Address
	data       CalculatorOperationData
}

// Init is called when the node is initialized by the simulation.
func (n *TestNode) Init(ctx context.Context) {
	n.data = CalculatorOperationData{A: n.A, B: n.B}
	n.SetTimer(ctx, ds.NewTimer(SendAdd, nil), 1*time.Second)
	n.SetTimer(ctx, ds.NewTimer(SendSubtract, nil), 2*time.Second)
	n.SetTimer(ctx, ds.NewTimer(SendMultiply, nil), 3*time.Second)
	n.SetTimer(ctx, ds.NewTimer(SendDivide, nil), 4*time.Second)
}

// HandleMessage is called when the node receives a message.
func (n *TestNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	case CalculatorResult:
		data := message.Data.(CalculatorResultData)
		fmt.Printf("Result: %v\n", data)
		return true
	default:
		return false
	}
}

// HandleTimer is called when a node receives a timer.
func (n *TestNode) HandleTimer(ctx context.Context, timer ds.Timer, duration time.Duration) bool {
	switch timer.Type {
	case SendAdd:
		n.SendMessage(ctx, ds.NewMessage(CalculatorAdd, n.data), n.calculator)
		n.SetTimer(ctx, ds.NewTimer(SendAdd, nil), 4*time.Second)
		return true
	case SendSubtract:
		n.SendMessage(ctx, ds.NewMessage(CalculatorSubtract, n.data), n.calculator)
		n.SetTimer(ctx, ds.NewTimer(SendSubtract, nil), 4*time.Second)
		return true
	case SendMultiply:
		n.SendMessage(ctx, ds.NewMessage(CalculatorMultiply, n.data), n.calculator)
		n.SetTimer(ctx, ds.NewTimer(SendMultiply, nil), 4*time.Second)
		return true
	case SendDivide:
		n.SendMessage(ctx, ds.NewMessage(CalculatorDivide, n.data), n.calculator)
		n.SetTimer(ctx, ds.NewTimer(SendDivide, nil), 4*time.Second)
		return true
	default:
		return false
	}
}
