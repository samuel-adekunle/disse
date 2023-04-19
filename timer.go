package disse

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TimerId is a string that uniquely identifies an interrupt in the network.
//
// It is generated using the github.com/google/uuid package.
type TimerId string

// TimerType is a string that identifies a timer type and is used to handle timers appropriately.
type TimerType string

// TimerData is the data associated with a timer.
type TimerData interface{}

// Timer is a message that is sent to a node after a certain amount of time to trigger certain events.
type Timer struct {
	Id   TimerId
	Type TimerType
	Data TimerData
}

// String returns a string representation of the timer for debugging purposes.
func (m Timer) String() string {
	return fmt.Sprintf("%v(%v, %v)", m.Type, m.Id, m.Data)
}

// NewTimer creates a new timer with the given timer type and data.
func NewTimer(timerType TimerType, data TimerData) Timer {
	return Timer{
		Id:   TimerId(uuid.NewString()),
		Type: timerType,
		Data: data,
	}
}

// TimerTriplet is a triplet of a timer, the address of the node to which it should be sent and the duration of the timer.
type TimerTriplet struct {
	Timer    Timer
	To       Address
	Duration time.Duration
}
