package disse

import (
	"log"
	"os"
	"time"
)

// Logger is an interface that is used to log events in the network.
//
// Each time an event occurs in the network, the corresponding Logger function is called.
type Logger interface {
	// Logger functions for state changes
	LogSimulationState(sim Simulation)
	LogNodeState(node Node)

	// Logger functions for messages
	LogSendMessage(from, to Address, message Message)
	LogHandleMessage(from, to Address, message Message)
	LogDropMessage(from, to Address, message Message)

	// Logger functions for timers
	LogSetTimer(to Address, timer Timer, duration time.Duration)
	LogHandleTimer(to Address, timer Timer, duration time.Duration)
	LogDropTimer(to Address, timer Timer, duration time.Duration)

	// Logger functions for interrupts
	LogSendInterrupt(from, to Address, interrupt Interrupt)
	LogHandleInterrupt(from, to Address, interrupt Interrupt)
	LogDropInterrupt(from, to Address, interrupt Interrupt)
}

// DebugLogger is a Log implementation that logs debug messages to a file.
type DebugLogger struct {
	logger *log.Logger
}

// NewDebugLogger creates a new DebugLog that logs to the given file.
func NewDebugLogger(logPath string) (*DebugLogger, error) {
	const (
		prefix = ""
		flag   = log.Ldate | log.Lmicroseconds
	)
	logfile, err := os.Create(logPath)
	if err != nil {
		return nil, err
	}
	return &DebugLogger{
		logger: log.New(logfile, prefix, flag),
	}, nil
}

// LogSimulationState is called when the simulation state changes.
func (l *DebugLogger) LogSimulationState(sim Simulation) {
	l.logger.Printf("SimulationState(%v)\n", sim.GetState())
}

// LogNodeState is called when the state of a node changes.
func (l *DebugLogger) LogNodeState(node Node) {
	l.logger.Printf("NodeState(%v, %v)\n", node.GetAddress(), node.GetState())
}

// LogSendMessage is called when a message is sent.
func (l *DebugLogger) LogSendMessage(from, to Address, message Message) {
	l.logger.Printf("SendMessage(%v -> %v, %v)\n", from, to, message)
}

// LogHandleMessage is called when a message is handled.
func (l *DebugLogger) LogHandleMessage(from, to Address, message Message) {
	l.logger.Printf("HandleMessage(%v -> %v, %v)\n", from, to, message)
}

// LogDropMessage is called when a message is dropped.
func (l *DebugLogger) LogDropMessage(from, to Address, message Message) {
	l.logger.Printf("DropMessage(%v -> %v, %v)\n", from, to, message)
}

// LogSetTimer is called when a timer is set.
func (l *DebugLogger) LogSetTimer(to Address, timer Timer, duration time.Duration) {
	l.logger.Printf("SetTimer(%v, %v, %v)\n", to, timer, duration)
}

// LogHandleTimer is called when a timer is handled.
func (l *DebugLogger) LogHandleTimer(to Address, timer Timer, duration time.Duration) {
	l.logger.Printf("HandleTimer(%v, %v, %v)\n", to, timer, duration)
}

// LogDropTimer is called when a timer is dropped.
func (l *DebugLogger) LogDropTimer(to Address, timer Timer, duration time.Duration) {
	l.logger.Printf("DropTimer(%v, %v, %v)\n", to, timer, duration)
}

// LogSendInterrupt is called when an interrupt is sent.
func (l *DebugLogger) LogSendInterrupt(from, to Address, interrupt Interrupt) {
	l.logger.Printf("SendInterrupt(%v -> %v, %v)\n", from, to, interrupt)
}

// LogHandleInterrupt is called when an interrupt is handled.
func (l *DebugLogger) LogHandleInterrupt(from, to Address, interrupt Interrupt) {
	l.logger.Printf("HandleInterrupt(%v -> %v, %v)\n", from, to, interrupt)
}

// LogDropInterrupt is called when an interrupt is dropped.
func (l *DebugLogger) LogDropInterrupt(from, to Address, interrupt Interrupt) {
	l.logger.Printf("DropInterrupt(%v -> %v, %v)\n", from, to, interrupt)
}

// UmlLogger is a Log implementation that logs messages in the PlantUML format.
type UmlLogger struct {
	logger *log.Logger
}

// NewUmlLogger creates a new UmlLog that logs to the given file.
func NewUmlLogger(logPath string) (*UmlLogger, error) {
	const (
		prefix = ""
		flag   = 0
	)
	logfile, err := os.Create(logPath)
	if err != nil {
		return nil, err
	}
	umlLog := log.New(logfile, prefix, flag)
	return &UmlLogger{
		logger: umlLog,
	}, nil
}

// LogSimulationState is called when the simulation state changes.
func (l *UmlLogger) LogSimulationState(sim Simulation) {
	switch sim.GetState() {
	case SimulationNotStarted:
		l.logger.Println("@startuml")
		l.logger.Println("!theme reddress-lightred")
		l.logger.Println("skinparam shadowing false")
		l.logger.Println("skinparam sequenceArrowThickness 1")
		l.logger.Println("skinparam responseMessageBelowArrow true")
		l.logger.Println("skinparam sequenceMessageAlign right")
	case SimulationRunning:
		// do nothing
	case SimulationFinished:
		l.logger.Println("@enduml")
	}
}

// LogNodeState is called when the state of a node changes.
func (l *UmlLogger) LogNodeState(node Node) {}

// LogSendMessage is called when a message is sent.
func (l *UmlLogger) LogSendMessage(from, to Address, message Message) {
	l.logger.Printf("%v -> %v : %v\n", from, to, message.Type)
}

// LogHandleMessage is called when a message is handled.
func (l *UmlLogger) LogHandleMessage(from, to Address, message Message) {}

// LogDropMessage is called when a message is dropped.
func (l *UmlLogger) LogDropMessage(from, to Address, message Message) {}

// LogSetTimer is called when a timer is set.
func (l *UmlLogger) LogSetTimer(to Address, timer Timer, duration time.Duration) {
	l.logger.Printf("%v -> %v : %v\n", to, to, timer.Type)
}

// LogHandleTimer is called when a timer is handled.
func (l *UmlLogger) LogHandleTimer(to Address, timer Timer, duration time.Duration) {}

// LogDropTimer is called when a timer is dropped.
func (l *UmlLogger) LogDropTimer(to Address, timer Timer, duration time.Duration) {}

// LogSendInterrupt is called when an interrupt is sent.
func (l *UmlLogger) LogSendInterrupt(from, to Address, interrupt Interrupt) {
	l.logger.Printf("%v -> %v : %v\n", from, to, interrupt.Type)
}

// LogHandleInterrupt is called when an interrupt is handled.
func (l *UmlLogger) LogHandleInterrupt(from, to Address, interrupt Interrupt) {}

// LogDropInterrupt is called when an interrupt is dropped.
func (l *UmlLogger) LogDropInterrupt(from, to Address, interrupt Interrupt) {}

// LogSimulationState is called when the simulation state changes.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogSimulationState() {
	for _, log := range s.loggers {
		log.LogSimulationState(s)
	}
}

// LogNodeState is called when the state of a node changes.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogNodeState(node Node) {
	for _, log := range s.loggers {
		log.LogNodeState(node)
	}
}

// LogSendMessage is called when a message is sent.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogSendMessage(from, to Address, message Message) {
	for _, log := range s.loggers {
		log.LogSendMessage(from, to, message)
	}
}

// LogHandleMessage is called when a message is handled.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogHandleMessage(from, to Address, message Message) {
	for _, log := range s.loggers {
		log.LogHandleMessage(from, to, message)
	}
}

// LogDropMessage is called when a message is dropped.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogDropMessage(from, to Address, message Message) {
	for _, log := range s.loggers {
		log.LogDropMessage(from, to, message)
	}
}

// LogSetTimer is called when a timer is set.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogSetTimer(to Address, timer Timer, duration time.Duration) {
	for _, log := range s.loggers {
		log.LogSetTimer(to, timer, duration)
	}
}

// LogHandleTimer is called when a timer is handled.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogHandleTimer(to Address, timer Timer, duration time.Duration) {
	for _, log := range s.loggers {
		log.LogHandleTimer(to, timer, duration)
	}
}

// LogDropTimer is called when a timer is dropped.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogDropTimer(to Address, timer Timer, duration time.Duration) {
	for _, log := range s.loggers {
		log.LogDropTimer(to, timer, duration)
	}
}

// LogSendInterrupt is called when an interrupt is sent.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogSendInterrupt(from, to Address, interrupt Interrupt) {
	for _, log := range s.loggers {
		log.LogSendInterrupt(from, to, interrupt)
	}
}

// LogHandleInterrupt is called when an interrupt is handled.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogHandleInterrupt(from, to Address, interrupt Interrupt) {
	for _, log := range s.loggers {
		log.LogHandleInterrupt(from, to, interrupt)
	}
}

// LogDropInterrupt is called when an interrupt is dropped.
//
// This method is called for all logs in the simulation.
func (s *LocalSimulation) LogDropInterrupt(from, to Address, interrupt Interrupt) {
	for _, log := range s.loggers {
		log.LogDropInterrupt(from, to, interrupt)
	}
}
