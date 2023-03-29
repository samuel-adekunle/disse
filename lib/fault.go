package lib

import (
	"time"
)

type FaultName string

const (
	Stop    FaultName = "StopFault"
	Resume  FaultName = "ResumeFault"
	Restart FaultName = "RestartFault"
	Sleep   FaultName = "SleepFault"
)

type Fault struct {
	Name     FaultName
	Duration time.Duration
}

type FaultTriplet struct {
	Fault Fault
	To    Address
	After time.Duration
}
