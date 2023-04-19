package disse

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Log is an interface that is used to log events in the network.
//
// Each time an event occurs in the network, the corresponding Log function is called.
type Log interface {
	// Log functions for state changes
	LogSimulationState(sim *Simulation)
	LogNodeState(node Node)

	// Log functions for messages
	LogSendMessage(from, to Address, message Message)
	LogHandleMessage(from, to Address, message Message)
	LogDropMessage(from, to Address, message Message)

	// Log functions for timers
	LogSetTimer(to Address, timer Timer, duration time.Duration)
	LogHandleTimer(to Address, timer Timer, duration time.Duration)
	LogDropTimer(to Address, timer Timer, duration time.Duration)

	// Log functions for interrupts
	LogSendInterrupt(from, to Address, interrupt Interrupt)
	LogHandleInterrupt(from, to Address, interrupt Interrupt)
	LogDropInterrupt(from, to Address, interrupt Interrupt)
}

// DebugLog is a Log implementation that logs debug messages to a file.
type DebugLog struct {
	log *log.Logger
}

// NewDebugLog creates a new DebugLog that logs to the given file.
func NewDebugLog(logPath string) *DebugLog {
	const (
		prefix = ""
		flag   = log.Ldate | log.Lmicroseconds
	)
	logfile, err := os.Create(logPath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	return &DebugLog{
		log: log.New(logfile, prefix, flag),
	}
}

// LogSimulationState is called when the simulation state changes.
func (l *DebugLog) LogSimulationState(sim *Simulation) {
	l.log.Printf("SimulationState(%v)\n", sim.state)
}

// LogNodeState is called when the state of a node changes.
func (l *DebugLog) LogNodeState(node Node) {
	l.log.Printf("NodeState(%v, %v)\n", node.GetAddress(), node.GetState())
}

// LogSendMessage is called when a message is sent.
func (l *DebugLog) LogSendMessage(from, to Address, message Message) {
	l.log.Printf("SendMessage(%v -> %v, %v)\n", from, to, message)
}

// LogHandleMessage is called when a message is handled.
func (l *DebugLog) LogHandleMessage(from, to Address, message Message) {
	l.log.Printf("HandleMessage(%v -> %v, %v)\n", from, to, message)
}

// LogDropMessage is called when a message is dropped.
func (l *DebugLog) LogDropMessage(from, to Address, message Message) {
	l.log.Printf("DropMessage(%v -> %v, %v)\n", from, to, message)
}

// LogSetTimer is called when a timer is set.
func (l *DebugLog) LogSetTimer(to Address, timer Timer, duration time.Duration) {
	l.log.Printf("SetTimer(%v, %v, %v)\n", to, timer, duration)
}

// LogHandleTimer is called when a timer is handled.
func (l *DebugLog) LogHandleTimer(to Address, timer Timer, duration time.Duration) {
	l.log.Printf("HandleTimer(%v, %v, %v)\n", to, timer, duration)
}

// LogDropTimer is called when a timer is dropped.
func (l *DebugLog) LogDropTimer(to Address, timer Timer, duration time.Duration) {
	l.log.Printf("DropTimer(%v, %v, %v)\n", to, timer, duration)
}

// LogSendInterrupt is called when an interrupt is sent.
func (l *DebugLog) LogSendInterrupt(from, to Address, interrupt Interrupt) {
	l.log.Printf("SendInterrupt(%v -> %v, %v)\n", from, to, interrupt)
}

// LogHandleInterrupt is called when an interrupt is handled.
func (l *DebugLog) LogHandleInterrupt(from, to Address, interrupt Interrupt) {
	l.log.Printf("HandleInterrupt(%v -> %v, %v)\n", from, to, interrupt)
}

// LogDropInterrupt is called when an interrupt is dropped.
func (l *DebugLog) LogDropInterrupt(from, to Address, interrupt Interrupt) {
	l.log.Printf("DropInterrupt(%v -> %v, %v)\n", from, to, interrupt)
}

// UmlLog is a Log implementation that logs messages in the PlantUML format.
type UmlLog struct {
	log *log.Logger
}

// NewUmlLog creates a new UmlLog that logs to the given file.
func NewUmlLog(logPath string) *UmlLog {
	const (
		prefix = ""
		flag   = 0
	)
	logfile, err := os.Create(logPath)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	}
	umlLog := log.New(logfile, prefix, flag)
	return &UmlLog{
		log: umlLog,
	}
}

// LogSimulationState is called when the simulation state changes.
func (l *UmlLog) LogSimulationState(sim *Simulation) {
	switch sim.state {
	case SimulationNotStarted:
		l.log.Println("@startuml")
		l.log.Println("!theme reddress-lightred")
		l.log.Println("skinparam shadowing false")
		l.log.Println("skinparam sequenceArrowThickness 1")
		l.log.Println("skinparam responseMessageBelowArrow true")
		l.log.Println("skinparam sequenceMessageAlign right")
	case SimulationRunning:
		// do nothing
	case SimulationFinished:
		l.log.Println("@enduml")
	}
}

// LogNodeState is called when the state of a node changes.
func (l *UmlLog) LogNodeState(node Node) {}

// LogSendMessage is called when a message is sent.
func (l *UmlLog) LogSendMessage(from, to Address, message Message) {
	l.log.Printf("%v -> %v : %v\n", from, to, message.Type)
}

// LogHandleMessage is called when a message is handled.
func (l *UmlLog) LogHandleMessage(from, to Address, message Message) {}

// LogDropMessage is called when a message is dropped.
func (l *UmlLog) LogDropMessage(from, to Address, message Message) {}

// LogSetTimer is called when a timer is set.
func (l *UmlLog) LogSetTimer(to Address, timer Timer, duration time.Duration) {
	l.log.Printf("%v -> %v : %v\n", to, to, timer.Type)
}

// LogHandleTimer is called when a timer is handled.
func (l *UmlLog) LogHandleTimer(to Address, timer Timer, duration time.Duration) {}

// LogDropTimer is called when a timer is dropped.
func (l *UmlLog) LogDropTimer(to Address, timer Timer, duration time.Duration) {}

// LogSendInterrupt is called when an interrupt is sent.
func (l *UmlLog) LogSendInterrupt(from, to Address, interrupt Interrupt) {
	l.log.Printf("%v -> %v : %v\n", from, to, interrupt.Type)
}

// LogHandleInterrupt is called when an interrupt is handled.
func (l *UmlLog) LogHandleInterrupt(from, to Address, interrupt Interrupt) {}

// LogDropInterrupt is called when an interrupt is dropped.
func (l *UmlLog) LogDropInterrupt(from, to Address, interrupt Interrupt) {}

// LogSimulationState is called when the simulation state changes.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogSimulationState() {
	for _, log := range s.loggers {
		log.LogSimulationState(s)
	}
}

// LogNodeState is called when the state of a node changes.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogNodeState(node Node) {
	for _, log := range s.loggers {
		log.LogNodeState(node)
	}
}

// LogSendMessage is called when a message is sent.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogSendMessage(from, to Address, message Message) {
	for _, log := range s.loggers {
		log.LogSendMessage(from, to, message)
	}
}

// LogHandleMessage is called when a message is handled.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogHandleMessage(from, to Address, message Message) {
	for _, log := range s.loggers {
		log.LogHandleMessage(from, to, message)
	}
}

// LogDropMessage is called when a message is dropped.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogDropMessage(from, to Address, message Message) {
	for _, log := range s.loggers {
		log.LogDropMessage(from, to, message)
	}
}

// LogSetTimer is called when a timer is set.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogSetTimer(to Address, timer Timer, duration time.Duration) {
	for _, log := range s.loggers {
		log.LogSetTimer(to, timer, duration)
	}
}

// LogHandleTimer is called when a timer is handled.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogHandleTimer(to Address, timer Timer, duration time.Duration) {
	for _, log := range s.loggers {
		log.LogHandleTimer(to, timer, duration)
	}
}

// LogDropTimer is called when a timer is dropped.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogDropTimer(to Address, timer Timer, duration time.Duration) {
	for _, log := range s.loggers {
		log.LogDropTimer(to, timer, duration)
	}
}

// LogSendInterrupt is called when an interrupt is sent.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogSendInterrupt(from, to Address, interrupt Interrupt) {
	for _, log := range s.loggers {
		log.LogSendInterrupt(from, to, interrupt)
	}
}

// LogHandleInterrupt is called when an interrupt is handled.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogHandleInterrupt(from, to Address, interrupt Interrupt) {
	for _, log := range s.loggers {
		log.LogHandleInterrupt(from, to, interrupt)
	}
}

// LogDropInterrupt is called when an interrupt is dropped.
//
// This method is called for all logs in the simulation.
func (s *Simulation) LogDropInterrupt(from, to Address, interrupt Interrupt) {
	for _, log := range s.loggers {
		log.LogDropInterrupt(from, to, interrupt)
	}
}
