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
type SimulationState string

const (
	// SimulationStateNotStarted is the state of the simulation before it is started.
	SimulationNotStarted SimulationState = "Not Started"
	// SimulationRunning is the state of the simulation while it is running.
	SimulationRunning SimulationState = "Running"
	// SimulationStateFinished is the state of the simulation after it is finished.
	SimulationFinished SimulationState = "Finished"
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
	loggers      []Log
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
		loggers: []Log{
			NewDebugLog(options.DebugLogPath),
			NewUmlLog(options.UmlLogPath),
		},
		state: SimulationNotStarted,
	}
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

// handleMessages handles a message once the appropriate node is found.
//
// If the node is not running, the message is dropped.
//
// If the message is successfully handled, true is returned, otherwise false.
func (s *Simulation) handleMessage(ctx context.Context, node Node, mt MessageTriplet) bool {
	if node.GetState() != Running {
		return false
	}
	s.LogHandleMessage(mt.From, mt.To, mt.Message)
	return node.HandleMessage(ctx, mt.Message, mt.From)
}

// HandleMessage handles a message by sending it to the appropriate node.
//
// If the root node does not exist, the message is dropped.
//
// If the node is not running, the message is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If the message is successfully handled, true is returned, otherwise false.
func (s *Simulation) HandleMessage(ctx context.Context, mt MessageTriplet) bool {
	if _, ok := s.nodes[mt.To.Root()]; !ok {
		return false
	}

	if mt.To == mt.To.Root() {
		return s.handleMessage(ctx, s.nodes[mt.To], mt)
	}

	return s.nodes[mt.To.Root()].FindMessageHandler(ctx, mt)
}

// DropMessage drops a message.
//
// This means the message is not handled by any node.
func (s *Simulation) DropMessage(ctx context.Context, mt MessageTriplet) {
	s.LogDropMessage(mt.From, mt.To, mt.Message)
}

// handleTimer handles a timer once the appropriate node is found.
//
// If the node is not running, the timer is dropped.
//
// If the timer is successfully handled, true is returned, otherwise false.
func (s *Simulation) handleTimer(ctx context.Context, node Node, tt TimerTriplet) bool {
	if node.GetState() != Running {
		return false
	}
	s.LogHandleTimer(tt.To, tt.Timer, tt.Duration)
	return node.HandleTimer(ctx, tt.Timer, tt.Duration)
}

// HandleTimer handles a timer by sending it to the appropriate node.
//
// If the node is not running, the timer is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If the timer is successfully handled, true is returned, otherwise false.
func (s *Simulation) HandleTimer(ctx context.Context, tt TimerTriplet) bool {
	if _, ok := s.nodes[tt.To.Root()]; !ok {
		return false
	}

	if tt.To == tt.To.Root() {
		return s.handleTimer(ctx, s.nodes[tt.To], tt)
	}

	return s.nodes[tt.To.Root()].FindTimerHandler(ctx, tt)
}

// DropTimer drops a timer.
//
// This means the timer is not handled by any node.
func (s *Simulation) DropTimer(ctx context.Context, tt TimerTriplet) {
	s.LogDropTimer(tt.To, tt.Timer, tt.Duration)
}

// handleInterrupt handles an interrupt once the appropriate node is found.
//
// If the node is not running, the interrupt is dropped.
//
// If the interrupt is successfully handled, true is returned, otherwise false.
func (s *Simulation) handleInterrupt(ctx context.Context, node Node, it InterruptTriplet) (handled bool) {
	if node.GetState() != Running {
		return false
	}
	s.LogHandleInterrupt(it.From, it.To, it.Interrupt)
	handled = node.HandleInterrupt(ctx, it.Interrupt, it.From)
	if handled {
		s.LogNodeState(node)
	}
	return handled
}

// HandleInterrupt handles an interrupt by sending it to the appropriate node.
//
// If the node is not running, the interrupt is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If an unknown interrupt is received, the interrupt is dropped and the function returns false, otherwise true.
func (s *Simulation) HandleInterrupt(ctx context.Context, it InterruptTriplet) bool {
	if _, ok := s.nodes[it.To.Root()]; !ok {
		return false
	}

	if it.To == it.To.Root() {
		return s.handleInterrupt(ctx, s.nodes[it.To], it)
	}

	return s.nodes[it.To.Root()].FindInterruptHandler(ctx, it)
}

// DropInterrupt drops an interrupt.
//
// This means the interrupt is not handled by any node.
func (s *Simulation) DropInterrupt(ctx context.Context, it InterruptTriplet) {
	s.LogDropInterrupt(it.From, it.To, it.Interrupt)
}

// randomLatency returns a random duration between the minimum and maximum latency.
func (s *Simulation) randomLatency() time.Duration {
	return s.options.MinLatency + time.Duration(rand.Int63n(int64(s.options.MaxLatency-s.options.MinLatency)))
}

// initNode initializes a node and all it's sub nodes.
func (s *Simulation) initNode(ctx context.Context, node Node) {
	node.Init(ctx)
	s.LogNodeState(node)
	node.InitAll(ctx)
}

// startSim starts the simulation by initializing all nodes and sub nodes.
func (s *Simulation) startSim(ctx context.Context) {
	s.LogSimulationState()
	for _, node := range s.nodes {
		s.initNode(ctx, node)
	}
	s.state = SimulationRunning
	s.LogSimulationState()
}

// stopSim stops the simulation by closing the message and timer queues and waiting for all nodes to stop doing work.
func (s *Simulation) stopSim() {
	close(s.messageQueue)
	close(s.timerQueue)
	s.state = SimulationFinished
	s.LogSimulationState()
}

// Run runs the simulation.
//
// The simulation run by polling the message and timer queues and sending the messages and timers to the appropriate nodes.
func (s *Simulation) Run() {
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
				s.generateUmlImage()
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

const (
	// JavaEnv is the environment variable name for the java executable.
	JavaEnv = "DISSE_JAVA"
	// PlantumlEnv is the environment variable name for the plantuml jar file.
	PlantumlEnv = "DISSE_PLANTUML"
)

// generateUmlImage generates a UML image of the simulation using PlantUML (requires java).
func (s *Simulation) generateUmlImage() {
	godotenv.Load()
	javaPath := os.Getenv(JavaEnv)
	plantumlPath := os.Getenv(PlantumlEnv)

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
