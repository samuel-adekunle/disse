package lib

import "time"

type Timer string

type TimerTriplet struct {
	Timer    Timer
	From     Address
	Duration time.Duration
}
