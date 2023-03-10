package lib

import "time"

type Timer string

type TimerTriplet struct {
	Timer    Timer
	To       Address
	Duration time.Duration
}
