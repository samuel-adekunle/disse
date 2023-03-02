package lib

import "time"

// Timer type, for all timers in the system
type Timer string

// type for timer	triplet timer, from, length
type TimerTriplet struct {
	Timer  Timer
	From   Address
	Length time.Duration
}
