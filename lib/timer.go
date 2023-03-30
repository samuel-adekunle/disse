package lib

import "time"

type TimerId string
type TimerData interface{}

type Timer struct {
	Id   TimerId
	Data TimerData
}

type TimerTriplet struct {
	Timer    Timer
	To       Address
	Duration time.Duration
}
