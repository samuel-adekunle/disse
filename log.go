package disse

import (
	"fmt"
	"log"
	"os"
	"time"
)

// TODO: call log functions in correct places in simulation.go and node.go without user having to call them

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
//
// The messages are user readable and contain the time, the file and the line number.
type DebugLog struct {
	log *log.Logger
}

// NewDebugLog creates a new DebugLog that logs to the given file.
func NewDebugLog(logPath string) *DebugLog {
	const (
		prefix = ""
		flag   = log.Ldate | log.Lmicroseconds | log.Lshortfile
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
// TODO: implement
func (l *DebugLog) LogSimulationState(sim *Simulation) {}

// LogNodeState is called when the state of a node changes.
// TODO: implement
func (l *DebugLog) LogNodeState(node Node) {}

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
//
// The messages are machine readable and can be used to generate UML images.
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
	umlLog.Println("@startuml")
	umlLog.Println("!theme reddress-lightred")
	umlLog.Println("skinparam shadowing false")
	umlLog.Println("skinparam sequenceArrowThickness 1")
	umlLog.Println("skinparam responseMessageBelowArrow true")
	umlLog.Println("skinparam sequenceMessageAlign right")
	return &UmlLog{
		log: umlLog,
	}
}

// LogSimulationState is called when the simulation state changes.
// TODO: implement
func (l *UmlLog) LogSimulationState(sim *Simulation) {
	// print when finished
	l.log.Println("@enduml")
}

// LogNodeState is called when the state of a node changes.
// TODO: implement
func (l *UmlLog) LogNodeState(node Node) {}

// LogSendMessage is called when a message is sent.
// TODO: implement
func (l *UmlLog) LogSendMessage(from, to Address, message Message) {}

// LogHandleMessage is called when a message is handled.
func (l *UmlLog) LogHandleMessage(from, to Address, message Message) {
	l.log.Printf("%v -> %v : %v\n", from, to, message)
}

// LogDropMessage is called when a message is dropped.
// TODO: implement
func (l *UmlLog) LogDropMessage(from, to Address, message Message) {}

// LogSetTimer is called when a timer is set.
// TODO: implement
func (l *UmlLog) LogSetTimer(to Address, timer Timer, duration time.Duration) {}

// LogHandleTimer is called when a timer is handled.
func (l *UmlLog) LogHandleTimer(to Address, timer Timer, duration time.Duration) {
	l.log.Printf("%v -> %v : %v\n", to, to, timer)
}

// LogDropTimer is called when a timer is dropped.
// TODO: implement
func (l *UmlLog) LogDropTimer(to Address, timer Timer, duration time.Duration) {}

// LogSendInterrupt is called when an interrupt is sent.
// TODO: implement
func (l *UmlLog) LogSendInterrupt(from, to Address, interrupt Interrupt) {}

// LogHandleInterrupt is called when an interrupt is handled.
func (l *UmlLog) LogHandleInterrupt(from, to Address, interrupt Interrupt) {
	l.log.Printf("%v -> %v : %v\n", from, to, interrupt)
}

// LogDropInterrupt is called when an interrupt is dropped.
// TODO: implement
func (l *UmlLog) LogDropInterrupt(from, to Address, interrupt Interrupt) {}
