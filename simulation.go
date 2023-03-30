package disse

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	Infinity                 = time.Duration(0)
	DefaultMessageBufferSize = 100
	DefaultTimerBufferSize   = 100
	DefaultMinLatency        = 10 * time.Millisecond
	DefaultMaxLatency        = 100 * time.Millisecond
	DefaultDuration          = Infinity
)

var logFileName string
var umlFileName string

// Init sets up the command line flags for the simulation executable.
// The log file name is the file where the simulation logs will be written.
// The default log file name is the name of the executable with the .log extension.
// The UML file name is the file where the UML diagram will be written.
// The default UML file name is the name of the executable with the .uml extension.
func init() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	baseFileName := strings.Split(workDir, "/")[len(strings.Split(workDir, "/"))-1]
	defaultLogFileName := fmt.Sprintf("%s.log", baseFileName)
	logFileNameUsage := "path to log file"
	defaultUmlFileName := fmt.Sprintf("%s.uml", baseFileName)
	umlFileNameUsage := "path to UML diagram file"

	flag.StringVar(&logFileName, "logfile", defaultLogFileName, "path to log file")
	flag.StringVar(&logFileName, "l", defaultLogFileName, logFileNameUsage+" (shorthand)")
	flag.StringVar(&umlFileName, "umlfile", defaultUmlFileName, "path to UML diagram file")
	flag.StringVar(&umlFileName, "u", defaultUmlFileName, umlFileNameUsage+" (shorthand)")
}

// Simulation sets up and runs the distributed system simulation.
type Simulation struct {
	nodes        map[Address]Node
	debugLog     *log.Logger
	umlLog       *log.Logger
	messageQueue chan MessageTriplet
	timerQueue   chan TimerTriplet
	MinLatency   time.Duration
	MaxLatency   time.Duration
	Duration     time.Duration
}

// BufferSizes is used to set the buffer sizes for the message and timer queues.
type BufferSizes struct {
	MessageBufferSize int
	TimerBufferSize   int
}

// NewSimulation creates a new simulation with the default buffer sizes.
func NewSimulation() *Simulation {
	return NewSimulationWithBuffer(nil)
}

// NewSimulationWithBuffer creates a new simulation with the given buffer sizes.
func NewSimulationWithBuffer(bufferSizes *BufferSizes) *Simulation {
	if bufferSizes == nil {
		bufferSizes = &BufferSizes{
			MessageBufferSize: DefaultMessageBufferSize,
			TimerBufferSize:   DefaultTimerBufferSize,
		}
	}
	return &Simulation{
		nodes:        make(map[Address]Node),
		messageQueue: make(chan MessageTriplet, bufferSizes.MessageBufferSize),
		timerQueue:   make(chan TimerTriplet, bufferSizes.TimerBufferSize),
		MinLatency:   DefaultMinLatency,
		MaxLatency:   DefaultMaxLatency,
		Duration:     DefaultDuration,
	}
}

// randomLatency returns a random duration between the minimum and maximum latency.
func (s *Simulation) randomLatency() time.Duration {
	return s.MinLatency + time.Duration(rand.Int63n(int64(s.MaxLatency-s.MinLatency)))
}

// AddNode adds a node to the simulation.
func (s *Simulation) AddNode(address Address, node Node) {
	s.nodes[address] = node
}

// AddNodes adds multiple nodes to the simulation.
//
// The addresses and nodes must be in the same order and have the same length if not an error is returned.
func (s *Simulation) AddNodes(addresses []Address, nodes []Node) (err error) {
	if len(addresses) != len(nodes) {
		return fmt.Errorf("length of addresses (%v) does not match length of nodes (%v)", len(addresses), len(nodes))
	}
	for i := range addresses {
		s.AddNode(addresses[i], nodes[i])
	}
	return nil
}

// HandleMessage handles a message by sending it to the appropriate node.
//
// If the node is not running, the message is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
func (s *Simulation) HandleMessage(ctx context.Context, mt MessageTriplet) {
	if node, ok := s.nodes[mt.To]; ok {
		if node.GetState() != Running {
			s.DropMessage(ctx, mt)
			return
		}
		s.debugLog.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message.Id)
		node.HandleMessage(ctx, mt.Message, mt.From)
	} else {
		s.nodes[mt.To.Root()].SubNodesHandleMessage(ctx, mt)
	}
}

// DropMessage drops a message.
func (s *Simulation) DropMessage(ctx context.Context, mt MessageTriplet) {
	s.debugLog.Printf("DropMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message.Id)
}

// HandleTimer handles a timer by sending it to the appropriate node.
//
// If the node is not running, the timer is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
func (s *Simulation) HandleTimer(ctx context.Context, tt TimerTriplet) {
	if node, ok := s.nodes[tt.To]; ok {
		if node.GetState() != Running {
			s.DropTimer(ctx, tt)
			return
		}
		s.debugLog.Printf("HandleTimer(%v, %v, %v)\n", tt.To, tt.Timer.Id, tt.Duration)
		node.HandleTimer(ctx, tt.Timer, tt.Duration)
	} else {
		s.nodes[tt.To.Root()].SubNodesHandleTimer(ctx, tt)
	}
}

// DropTimer drops a timer.
func (s *Simulation) DropTimer(ctx context.Context, tt TimerTriplet) {
	s.debugLog.Printf("DropTimer(%v, %v, %v)\n", tt.To, tt.Timer.Id, tt.Duration)
}

// HandleInterrupt handles an interrupt by sending it to the appropriate node.
//
// If the node is not running, the interrupt is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If an unknown interrupt is received, an error is returned.
func (s *Simulation) HandleInterrupt(ctx context.Context, ip InterruptPair) (err error) {
	if node, ok := s.nodes[ip.To]; ok {
		if node.GetState() == Stopped {
			s.DropInterrupt(ctx, ip)
			return
		}
		s.debugLog.Printf("HandleInterrupt(%v, %v)\n", ip.To, ip.Interrupt.Id)
		err = node.HandleInterrupt(ctx, ip.Interrupt)
	} else {
		err = s.nodes[ip.To.Root()].SubNodesHandleInterrupt(ctx, ip)
	}
	return
}

// DropInterrupt drops an interrupt.
func (s *Simulation) DropInterrupt(ctx context.Context, ip InterruptPair) {
	s.debugLog.Printf("DropInterrupt(%v, %v)\n", ip.To, ip.Interrupt.Id)
}

// startSim starts the simulation by initializing all nodes and sub nodes.
func (s *Simulation) startSim(ctx context.Context) {
	s.debugLog.Printf("StartSim(%v)\n", s.Duration)
	for address, node := range s.nodes {
		s.debugLog.Printf("Init(%v)\n", address)
		node.Init(ctx)
		node.SubNodesInit(ctx)
	}
}

// stopSim stops the simulation by closing the message and timer queues and waiting for all nodes to stop doing work.
func (s *Simulation) stopSim() {
	s.debugLog.Println("StopSim()")
	close(s.messageQueue)
	close(s.timerQueue)
}

// Run runs the simulation.
//
// The simulation run by polling the message and timer queues and sending the messages and timers to the appropriate nodes.
func (s *Simulation) Run() {
	flag.Parse()
	logFile, err := os.Create(logFileName)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	s.debugLog = log.New(logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	umlFile, err := os.Create(umlFileName)
	if err != nil {
		panic(err)
	}
	defer umlFile.Close()
	s.umlLog = log.New(umlFile, "", 0)
	s.umlLog.Println("@startuml")
	s.umlLog.Println("!theme reddress-lightred")
	s.umlLog.Println("skinparam shadowing false")
	s.umlLog.Println("skinparam sequenceArrowThickness 1")
	s.umlLog.Println("skinparam responseMessageBelowArrow true")
	s.umlLog.Println("skinparam sequenceMessageAlign right")
	defer s.umlLog.Println("@enduml")

	var ctx context.Context
	if s.Duration == Infinity {
		ctx = context.Background()
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.Duration)
		defer cancel()
	}

	var wg sync.WaitGroup
	s.startSim(ctx)
	for {
		select {
		case <-ctx.Done():
			if s.Duration != Infinity {
				s.stopSim()
				wg.Wait()
				s.debugLog.Println("EndSim()")
				return
			}
		case mt := <-s.messageQueue:
			wg.Add(1)
			go func() {
				time.Sleep(s.randomLatency())
				s.HandleMessage(ctx, mt)
				wg.Done()
			}()
		case tt := <-s.timerQueue:
			wg.Add(1)
			go func() {
				time.Sleep(tt.Duration)
				s.HandleTimer(ctx, tt)
				wg.Done()
			}()
		}
	}
}
