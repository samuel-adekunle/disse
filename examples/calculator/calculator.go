package main

import (
	"context"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const (
	// CalculatorAdd is the type of message used to add two numbers.
	CalculatorAdd ds.MessageType = "CalculatorAdd"
	// CalculatorSubtract is the type of message used to subtract two numbers.
	CalculatorSubtract ds.MessageType = "CalculatorSubtract"
	// CalculatorMultiply is the type of message used to multiply two numbers.
	CalculatorMultiply ds.MessageType = "CalculatorMultiply"
	// CalculatorDivide is the type of message used to divide two numbers.
	CalculatorDivide ds.MessageType = "CalculatorDivide"
	// CalculatorResult is the type of message used to send the result of an operation.
	CalculatorResult ds.MessageType = "CalculatorResult"
)

// CalculatorOperationData is the data of a CalculatorAdd,
// CalculatorSubtract, CalculatorMultiply or CalculatorDivide message.
type CalculatorOperationData struct {
	A int
	B int
}

// CalculatorResultData is the data of a CalculatorResult message.
type CalculatorResultData struct {
	A         int
	B         int
	Operation ds.MessageType
	Result    int
}

// CalculatorNode is a node that performs calculations.
type CalculatorNode struct {
	*ds.LocalNode
}

// Init is called when the node is initialized by the simulation.
func (n *CalculatorNode) Init(ctx context.Context) {}

// HandleMessage is called when the node receives a message.
func (n *CalculatorNode) HandleMessage(ctx context.Context, message ds.Message, from ds.Address) bool {
	switch message.Type {
	default:
		return false
	}
}

// HandleTimer is called when a node receives a timer.
// HandleTimer(context.Context, Timer, time.Duration) (handled bool)
func (n *CalculatorNode) HandleTimer(ctx context.Context, timer ds.Timer, duration time.Duration) bool {
	switch timer.Type {
	default:
		return false
	}
}
