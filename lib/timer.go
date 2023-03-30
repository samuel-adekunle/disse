package lib

import "time"

type TimerId string
type TimerData interface{}

type Timer struct {
	Id   TimerId
	Data TimerData
}

func NewTimer(id TimerId, data TimerData) Timer {
	return Timer{
		Id:   id,
		Data: data,
	}
}

type TimerTriplet struct {
	Timer    Timer
	To       Address
	Duration time.Duration
}
