package disse

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// InterruptId is a string that uniquely identifies an interrupt in the network.
//
// It is generated using the github.com/google/uuid package.
type InterruptId string

// InterruptType is a string that identifies an interrupt type.
type InterruptType string

const (
	StopInterrupt  InterruptType = "StopInterrupt"
	StartInterrupt InterruptType = "StartInterrupt"
	SleepInterrupt InterruptType = "SleepInterrupt"
)

// InterruptData is the data associated with an interrupt.
type InterruptData interface{}

// Interrupt is a message that is sent to a node to interrupt its execution in some way.
//
// It is used to stop a node, to make it sleep for a while or to make it start again.
type Interrupt struct {
	Id   InterruptId
	Type InterruptType
	Data InterruptData
}

// String returns a string representation of the interrupt for debugging purposes.
func (i Interrupt) String() string {
	return fmt.Sprintf("%v(%v, %v)", i.Type, i.Id, i.Data)
}

// NewInterrupt creates a new interrupt with the given interrupt type and data.
func NewInterrupt(interruptType InterruptType, data InterruptData) Interrupt {
	return Interrupt{
		Id:   InterruptId(uuid.NewString()),
		Type: interruptType,
		Data: data,
	}
}

// SleepInterruptData is the data associated with a SleepInterrupt.
type SleepInterruptData struct {
	Duration time.Duration
}

// InterruptTriplet is a triplet of an interrupt, the address of the node that sent the interrupt and the address of the node that received the interrupt.
type InterruptTriplet struct {
	Interrupt Interrupt
	From      Address
	To        Address
}
