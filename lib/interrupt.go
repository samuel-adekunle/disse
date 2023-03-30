package lib

import "time"

type InterruptId string
type InterruptData interface{}

const (
	StopInterrupt  InterruptId = "StopInterrupt"
	StartInterrupt InterruptId = "StartInterrupt"
	SleepInterrupt InterruptId = "SleepInterrupt"
)

type Interrupt struct {
	Id   InterruptId
	Data InterruptData
}

type SleepInterruptData struct {
	Duration time.Duration
}

type InterruptPair struct {
	Interrupt Interrupt
	To        Address
}
