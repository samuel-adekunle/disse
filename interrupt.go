package disse

import "time"

// InterruptId is a string that identifies an interrupt type.
type InterruptId string

const (
	StopInterrupt  InterruptId = "StopInterrupt"
	StartInterrupt InterruptId = "StartInterrupt"
	SleepInterrupt InterruptId = "SleepInterrupt"
)

// InterruptData is the data associated with an interrupt.
type InterruptData interface{}

// Interrupt is a message that is sent to a node to interrupt its execution in some way.
//
// It is used to stop a node, to make it sleep for a while or to make it start again.
type Interrupt struct {
	Id   InterruptId
	Data InterruptData
}

// SleepInterruptData is the data associated with a SleepInterrupt.
type SleepInterruptData struct {
	Duration time.Duration
}

// InterruptPair is a pair of an interrupt and the address of the node to which it should be sent.
type InterruptPair struct {
	Interrupt Interrupt
	To        Address
}
