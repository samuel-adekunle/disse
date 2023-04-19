package disse

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// SimulationState is the state of the simulation.
type SimulationState int

const (
	// SimulationStateNotStarted is the state of the simulation before it is started.
	SimulationNotStarted SimulationState = iota
	// SimulationRunning is the state of the simulation while it is running.
	SimulationRunning
	// SimulationStateFinished is the state of the simulation after it is finished.
	SimulationFinished
)

// SimulationOptions is used to set the options for the simulation.
//
//   - MinLatency is the minimum latency of the network.
//   - MaxLatency is the maximum latency of the network.
//   - Duration is the duration of the simulation. If it is set to Infinity, the simulation will run forever.
//   - MessageBufferSize is the size of the message queue.
//   - TimerBufferSize is the size of the timer queue.
//   - DebugLogPath is the path to the debug log file.
//   - UmlLogPath is the path to the UML log file.
type SimulationOptions struct {
	MinLatency        time.Duration
	MaxLatency        time.Duration
	Duration          time.Duration
	MessageBufferSize int
	TimerBufferSize   int
	DebugLogPath      string
	UmlLogPath        string
}

const (
	Infinity                 = time.Duration(0)
	DefaultMessageBufferSize = 100
	DefaultTimerBufferSize   = 100
	DefaultMinLatency        = 10 * time.Millisecond
	DefaultMaxLatency        = 100 * time.Millisecond
	DefaultDuration          = Infinity
	DefaultDebugLogPath      = ""
	DefaultUmlLogPath        = ""
)

// Simulation sets up and runs the distributed system simulation.
type Simulation struct {
	options      *SimulationOptions
	nodes        map[Address]Node
	messageQueue chan MessageTriplet
	timerQueue   chan TimerTriplet
	debugLog     Log
	umlLog       Log
	state        SimulationState
}

// NewSimulation creates a new simulation with the given options.
//
// If the options are nil, the default options are used.
//
// The default options are:
//   - MinLatency: 10ms
//   - MaxLatency: 100ms
//   - Duration: Infinity (runs forever)
//   - MessageBufferSize: 100
//   - TimerBufferSize: 100
//   - DebugLogPath: "" (no debug log)
//   - UmlLogPath: "" (no UML log)
func NewSimulation(options *SimulationOptions) *Simulation {
	if options == nil {
		options = &SimulationOptions{
			MinLatency:        DefaultMinLatency,
			MaxLatency:        DefaultMaxLatency,
			Duration:          DefaultDuration,
			MessageBufferSize: DefaultMessageBufferSize,
			TimerBufferSize:   DefaultTimerBufferSize,
			DebugLogPath:      DefaultDebugLogPath,
			UmlLogPath:        DefaultUmlLogPath,
		}
	}
	return &Simulation{
		options:      options,
		nodes:        make(map[Address]Node),
		messageQueue: make(chan MessageTriplet, options.MessageBufferSize),
		timerQueue:   make(chan TimerTriplet, options.TimerBufferSize),
		debugLog:     NewDebugLog(options.DebugLogPath),
		umlLog:       NewUmlLog(options.UmlLogPath),
		state:        SimulationNotStarted,
	}
}

// randomLatency returns a random duration between the minimum and maximum latency.
func (s *Simulation) randomLatency() time.Duration {
	return s.options.MinLatency + time.Duration(rand.Int63n(int64(s.options.MaxLatency-s.options.MinLatency)))
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
//
// If the message is successfully handled, true is returned, otherwise false.
func (s *Simulation) HandleMessage(ctx context.Context, mt MessageTriplet) (handled bool) {
	if node, ok := s.nodes[mt.To]; ok {
		if node.GetState() != Running {
			return false
		}
		s.debugLog.LogHandleMessage(mt.From, mt.To, mt.Message)
		return node.HandleMessage(ctx, mt.Message, mt.From)
	}
	return s.nodes[mt.To.Root()].SubNodesHandleMessage(ctx, mt)
}

// DropMessage drops a message.
//
// This means the message is not handled by any node.
func (s *Simulation) DropMessage(ctx context.Context, mt MessageTriplet) {
	s.debugLog.LogDropMessage(mt.From, mt.To, mt.Message)
}

// HandleTimer handles a timer by sending it to the appropriate node.
//
// If the node is not running, the timer is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If the timer is successfully handled, true is returned, otherwise false.
func (s *Simulation) HandleTimer(ctx context.Context, tt TimerTriplet) (handled bool) {
	if node, ok := s.nodes[tt.To]; ok {
		if node.GetState() != Running {
			return false
		}
		s.debugLog.LogHandleTimer(tt.To, tt.Timer, tt.Duration)
		return node.HandleTimer(ctx, tt.Timer, tt.Duration)
	}
	return s.nodes[tt.To.Root()].SubNodesHandleTimer(ctx, tt)
}

// DropTimer drops a timer.
//
// This means the timer is not handled by any node.
func (s *Simulation) DropTimer(ctx context.Context, tt TimerTriplet) {
	s.debugLog.LogDropTimer(tt.To, tt.Timer, tt.Duration)
}

// HandleInterrupt handles an interrupt by sending it to the appropriate node.
//
// If the node is not running, the interrupt is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If an unknown interrupt is received, the interrupt is dropped and the function returns false, otherwise true.
func (s *Simulation) HandleInterrupt(ctx context.Context, ip InterruptTriplet) bool {
	if node, ok := s.nodes[ip.To]; ok {
		if node.GetState() == Stopped {
			return false
		}
		s.debugLog.LogHandleInterrupt(ip.From, ip.To, ip.Interrupt)
		return node.HandleInterrupt(ctx, ip.Interrupt)
	}
	return s.nodes[ip.To.Root()].SubNodesHandleInterrupt(ctx, ip)
}

// DropInterrupt drops an interrupt.
//
// This means the interrupt is not handled by any node.
func (s *Simulation) DropInterrupt(ctx context.Context, ip InterruptTriplet) {
	s.debugLog.LogDropInterrupt(ip.From, ip.To, ip.Interrupt)
}

// startSim starts the simulation by initializing all nodes and sub nodes.
func (s *Simulation) startSim(ctx context.Context) {
	s.debugLog.LogSimulationState(s)
	for _, node := range s.nodes {
		s.debugLog.LogNodeState(node)
		node.Init(ctx)
		node.SubNodesInit(ctx)
	}
}

// stopSim stops the simulation by closing the message and timer queues and waiting for all nodes to stop doing work.
func (s *Simulation) stopSim() {
	s.debugLog.LogSimulationState(s)
	close(s.messageQueue)
	close(s.timerQueue)
	s.generateUmlImage()
}

// Run runs the simulation.
//
// The simulation run by polling the message and timer queues and sending the messages and timers to the appropriate nodes.
func (s *Simulation) Run() {
	godotenv.Load()
	var ctx context.Context
	if s.options.Duration == Infinity {
		ctx = context.Background()
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.options.Duration)
		defer cancel()
	}

	var wg sync.WaitGroup
	s.startSim(ctx)
	for {
		select {
		case <-ctx.Done():
			if s.options.Duration != Infinity {
				s.stopSim()
				wg.Wait()
				s.debugLog.LogSimulationState(s)
				return
			}
		case mt := <-s.messageQueue:
			wg.Add(1)
			go func() {
				time.Sleep(s.randomLatency())
				if handled := s.HandleMessage(ctx, mt); !handled {
					s.DropMessage(ctx, mt)
				}
				wg.Done()
			}()
		case tt := <-s.timerQueue:
			wg.Add(1)
			go func() {
				time.Sleep(tt.Duration)
				if handled := s.HandleTimer(ctx, tt); !handled {
					s.DropTimer(ctx, tt)
				}
				wg.Done()
			}()
		}
	}
}

var (
	javaEnv     = "DISSE_JAVA"
	plantumlEnv = "DISSE_PLANTUML"
)

// generateUmlImage generates a UML image of the simulation using PlantUML (requires java).
func (s *Simulation) generateUmlImage() {
	javaPath := os.Getenv(javaEnv)
	plantumlPath := os.Getenv(plantumlEnv)

	if javaPath == "" || plantumlPath == "" {
		fmt.Printf("javaPath (%v) or plantumlPath (%v) not set. UML image not generated.\n", javaPath, plantumlPath)
		return
	}

	if s.options.UmlLogPath == "" {
		fmt.Printf("umlLogPath (%v) set to ''. UML image not generated.\n", s.options.UmlLogPath)
		return
	}

	cmd := exec.Command(javaPath, "-jar", plantumlPath, s.options.UmlLogPath)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error %v when generating UML image. Check if javaPath (%v) and plantumlPath (%v) are correctly set.\n", err, javaPath, plantumlPath)
	}
}
