package disse

import (
	"context"
	"fmt"
	"math/rand"
	"os/exec"
	"sync"
	"time"
)

// Simulation is the interface that must be implemented by all simulations.
type Simulation interface {
	AddNode(Node)
	RemoveNode(Address)
	AddLogger(Logger)
	RemoveLogger(Logger)
	Run()
}

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
type SimulationOptions struct {
	MinLatency   time.Duration
	MaxLatency   time.Duration
	Duration     time.Duration
	BufferSize   int
	DebugLogPath string
	UmlLogPath   string
	JavaPath     string
	PlantumlPath string
}

const (
	// DefaultBufferSize is the default buffer size for the message queue.
	DefaultBufferSize = 20
	// DefaultMinLatency is the default minimum latency for messages.
	DefaultMinLatency = 10 * time.Millisecond
	// DefaultMaxLatency is the default maximum latency for messages.
	DefaultMaxLatency = 100 * time.Millisecond
	// DefaultDuration is the default duration of the simulation.
	DefaultDuration = 10 * time.Second
	// DefaultDebugLogPath is the default path to the debug log.
	DefaultDebugLogPath = "debug.log"
	// DefaultUmlLogPath is the default path to the UML log.
	DefaultUmlLogPath = "uml.log"
	// DefaultJavaPath is the default path to the java executable.
	DefaultJavaPath = ""
	// DefaultPlantumlPath is the default path to the plantuml jar file.
	DefaultPlantumlPath = ""
)

// LocalSimulation sets up and runs the distributed system simulation locally using shared memory.
type LocalSimulation struct {
	options        *SimulationOptions
	nodes          map[Address]Node
	wg             *sync.WaitGroup
	messageQueue   map[Address]chan MessageTriplet
	timerQueue     map[Address]chan TimerTriplet
	interruptQueue map[Address]chan InterruptTriplet
	loggers        []Logger
	state          SimulationState
}

// NewLocalSimulation creates a new simulation with the given options.
//
// If the options are nil, the default options are used.
func NewLocalSimulation(options *SimulationOptions) *LocalSimulation {
	if options == nil {
		options = &SimulationOptions{
			MinLatency:   DefaultMinLatency,
			MaxLatency:   DefaultMaxLatency,
			Duration:     DefaultDuration,
			BufferSize:   DefaultBufferSize,
			DebugLogPath: DefaultDebugLogPath,
			UmlLogPath:   DefaultUmlLogPath,
			JavaPath:     DefaultJavaPath,
			PlantumlPath: DefaultPlantumlPath,
		}
	}
	return &LocalSimulation{
		options:        options,
		wg:             &sync.WaitGroup{},
		nodes:          make(map[Address]Node),
		messageQueue:   make(map[Address]chan MessageTriplet),
		timerQueue:     make(map[Address]chan TimerTriplet),
		interruptQueue: make(map[Address]chan InterruptTriplet),
		loggers:        make([]Logger, 0),
		state:          SimulationNotStarted,
	}
}

// AddNode adds a node to the simulation.
func (s *LocalSimulation) AddNode(node Node) {
	address := node.GetAddress()
	s.nodes[address] = node
	s.messageQueue[address] = make(chan MessageTriplet, s.options.BufferSize)
	s.timerQueue[address] = make(chan TimerTriplet, s.options.BufferSize)
	s.interruptQueue[address] = make(chan InterruptTriplet, s.options.BufferSize)
}

// RemoveNode removes a node from the simulation.
func (s *LocalSimulation) RemoveNode(address Address) {
	delete(s.nodes, address)
	delete(s.messageQueue, address)
	delete(s.timerQueue, address)
	delete(s.interruptQueue, address)
}

// AddLogger adds a logger to the simulation.
func (s *LocalSimulation) AddLogger(logger Logger) {
	s.loggers = append(s.loggers, logger)
}

// RemoveLogger removes a logger from the simulation.
func (s *LocalSimulation) RemoveLogger(logger Logger) {
	for i, l := range s.loggers {
		if l == logger {
			s.loggers = append(s.loggers[:i], s.loggers[i+1:]...)
			return
		}
	}
}

// handleMessages handles a message once the appropriate node is found.
//
// If the node is not running, the message is dropped.
//
// If the message is successfully handled, true is returned, otherwise false.
func (s *LocalSimulation) _handleMessage(ctx context.Context, node Node, mt MessageTriplet) bool {
	if node.GetState() != Running {
		return false
	}
	s.LogHandleMessage(mt.From, mt.To, mt.Message)
	return node.HandleMessage(ctx, mt.Message, mt.From)
}

// handleMessage handles a message by sending it to the appropriate node.
//
// If the root node does not exist, the message is dropped.
//
// If the node is not running, the message is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If the message is successfully handled, true is returned, otherwise false.
func (s *LocalSimulation) handleMessage(ctx context.Context, mt MessageTriplet) bool {
	if _, ok := s.nodes[mt.To.Root()]; !ok {
		return false
	}

	if mt.To == mt.To.Root() {
		return s._handleMessage(ctx, s.nodes[mt.To], mt)
	}

	return s.nodes[mt.To.Root()].findMessageHandler(ctx, mt)
}

// dropMessage drops a message.
//
// This means the message is not handled by any node.
func (s *LocalSimulation) dropMessage(ctx context.Context, mt MessageTriplet) {
	s.LogDropMessage(mt.From, mt.To, mt.Message)
}

// _handleTimer handles a timer once the appropriate node is found.
//
// If the node is not running, the timer is dropped.
//
// If the timer is successfully handled, true is returned, otherwise false.
func (s *LocalSimulation) _handleTimer(ctx context.Context, node Node, tt TimerTriplet) bool {
	if node.GetState() != Running {
		return false
	}
	s.LogHandleTimer(tt.To, tt.Timer, tt.Duration)
	return node.HandleTimer(ctx, tt.Timer, tt.Duration)
}

// handleTimer handles a timer by sending it to the appropriate node.
//
// If the node is not running, the timer is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If the timer is successfully handled, true is returned, otherwise false.
func (s *LocalSimulation) handleTimer(ctx context.Context, tt TimerTriplet) bool {
	if _, ok := s.nodes[tt.To.Root()]; !ok {
		return false
	}

	if tt.To == tt.To.Root() {
		return s._handleTimer(ctx, s.nodes[tt.To], tt)
	}

	return s.nodes[tt.To.Root()].findTimerHandler(ctx, tt)
}

// dropTimer drops a timer.
//
// This means the timer is not handled by any node.
func (s *LocalSimulation) dropTimer(ctx context.Context, tt TimerTriplet) {
	s.LogDropTimer(tt.To, tt.Timer, tt.Duration)
}

// _handleInterrupt handles an interrupt once the appropriate node is found.
//
// If the node is not running, the interrupt is dropped.
//
// If the interrupt is successfully handled, true is returned, otherwise false.
func (s *LocalSimulation) _handleInterrupt(ctx context.Context, node Node, it InterruptTriplet) (handled bool) {
	if node.GetState() != Running {
		return false
	}
	s.LogHandleInterrupt(it.From, it.To, it.Interrupt)
	handled = node.handleInterrupt(ctx, it.Interrupt, it.From)
	if handled {
		s.LogNodeState(node)
	}
	return handled
}

// handleInterrupt handles an interrupt by sending it to the appropriate node.
//
// If the node is not running, the interrupt is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If an unknown interrupt is received, the interrupt is dropped and the function returns false, otherwise true.
func (s *LocalSimulation) handleInterrupt(ctx context.Context, it InterruptTriplet) bool {
	if _, ok := s.nodes[it.To.Root()]; !ok {
		return false
	}

	if it.To == it.To.Root() {
		return s._handleInterrupt(ctx, s.nodes[it.To], it)
	}

	return s.nodes[it.To.Root()].findInterruptHandler(ctx, it)
}

// dropInterrupt drops an interrupt.
//
// This means the interrupt is not handled by any node.
func (s *LocalSimulation) dropInterrupt(ctx context.Context, it InterruptTriplet) {
	s.LogDropInterrupt(it.From, it.To, it.Interrupt)
}

// randomLatency returns a random duration between the minimum and maximum latency.
func (s *LocalSimulation) randomLatency() time.Duration {
	return s.options.MinLatency + time.Duration(rand.Int63n(int64(s.options.MaxLatency-s.options.MinLatency)))
}

// initNode initializes a node and all it's sub nodes.
func (s *LocalSimulation) initNode(ctx context.Context, node Node) {
	node.Init(ctx)
	s.LogNodeState(node)
	node.initSubNodes(ctx)
}

// startSim starts the simulation by initializing all nodes and sub nodes.
func (s *LocalSimulation) startSim(ctx context.Context) {
	s.LogSimulationState()
	for _, node := range s.nodes {
		s.initNode(ctx, node)
		s.wg.Add(1)
		go func(_node Address) {
			for {
				select {
				case <-ctx.Done():
					s.wg.Done()
					return
				case mt := <-s.messageQueue[_node]:
					if handled := s.handleMessage(ctx, mt); !handled {
						s.dropMessage(ctx, mt)
					}
				case tt := <-s.timerQueue[_node]:
					if handled := s.handleTimer(ctx, tt); !handled {
						s.dropTimer(ctx, tt)
					}
				case it := <-s.interruptQueue[_node]:
					if handled := s.handleInterrupt(ctx, it); !handled {
						s.dropInterrupt(ctx, it)
					}
				}
			}
		}(node.GetAddress())
	}
	s.state = SimulationRunning
	s.LogSimulationState()
}

// stopSim stops the simulation by closing the message and timer queues and waiting for all nodes to stop doing work.
func (s *LocalSimulation) stopSim() {
	s.wg.Wait()
	s.state = SimulationFinished
	s.LogSimulationState()
}

// Run runs the simulation.
func (s *LocalSimulation) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), s.options.Duration)
	defer cancel()

	debugLog := NewDebugLog(s.options.DebugLogPath)
	if debugLog != nil {
		s.AddLogger(debugLog)
	}

	umlLog := NewUmlLog(s.options.UmlLogPath)
	if umlLog != nil {
		s.AddLogger(umlLog)
	}

	s.startSim(ctx)
	<-ctx.Done()
	s.stopSim()
	s.generateUmlImage()
}

// generateUmlImage generates a UML image of the simulation using PlantUML (requires java).
func (s *LocalSimulation) generateUmlImage() {
	javaPath := s.options.JavaPath
	plantumlPath := s.options.PlantumlPath
	if javaPath == "" || plantumlPath == "" {
		fmt.Println("javaPath or plantumlPath not set. UML image not generated.")
		return
	}

	if s.options.UmlLogPath == "" {
		fmt.Println("umlLogPath not set. UML image not generated.")
		return
	}

	cmd := exec.Command(javaPath, "-jar", plantumlPath, s.options.UmlLogPath)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error %v when generating UML image. Check if javaPath, plantumlPath and umlLogPath.\n", err)
	}
}
