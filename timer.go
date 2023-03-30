package disse

import "time"

// TimerId is a string that identifies a timer type and is used to handle timers appropriately.
type TimerId string

// TimerData is the data associated with a timer.
type TimerData interface{}

// Timer is a message that is sent to a node after a certain amount of time to trigger certain events.
type Timer struct {
	Id   TimerId
	Data TimerData
}

// NewTimer creates a new timer with the given id and data.
func NewTimer(id TimerId, data TimerData) Timer {
	return Timer{
		Id:   id,
		Data: data,
	}
}

// TimerTriplet is a triplet of a timer, the address of the node to which it should be sent and the duration of the timer.
type TimerTriplet struct {
	Timer    Timer
	To       Address
	Duration time.Duration
}
