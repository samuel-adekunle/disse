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
	GetState() SimulationState
	GetOptions() LocalSimulationOptions
	GetNodes() map[Address]Node
	AddNode(Node) error
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

// LocalSimulationOptions is used to set the options for the simulation.
type LocalSimulationOptions struct {
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
	DefaultBufferSize = 10
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
	DefaultJavaPath = "/usr/bin/java"
	// DefaultPlantumlPath is the default path to the plantuml jar file.
	DefaultPlantumlPath = "/usr/share/plantuml/plantuml.jar"
)

// LocalSimulation sets up and runs the distributed system simulation locally using shared memory.
type LocalSimulation struct {
	options        *LocalSimulationOptions
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
func NewLocalSimulation(options *LocalSimulationOptions) *LocalSimulation {
	if options == nil {
		options = &LocalSimulationOptions{
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
	sim := &LocalSimulation{
		options:        options,
		wg:             &sync.WaitGroup{},
		nodes:          make(map[Address]Node),
		messageQueue:   make(map[Address]chan MessageTriplet),
		timerQueue:     make(map[Address]chan TimerTriplet),
		interruptQueue: make(map[Address]chan InterruptTriplet),
		loggers:        make([]Logger, 0),
		state:          SimulationNotStarted,
	}
	debugLogger, _ := NewDebugLogger(options.DebugLogPath)
	sim.AddLogger(debugLogger)
	umlLogger, _ := NewUmlLogger(options.UmlLogPath)
	sim.AddLogger(umlLogger)
	return sim
}

// GetState returns the state of the simulation.
func (s *LocalSimulation) GetState() SimulationState {
	return s.state
}

// GetOptions returns the options of the simulation.
func (s *LocalSimulation) GetOptions() LocalSimulationOptions {
	return *s.options
}

// GetNodes returns the nodes in the simulation.
func (s *LocalSimulation) GetNodes() map[Address]Node {
	return s.nodes
}

// AddNode adds a node to the simulation.
func (s *LocalSimulation) AddNode(node Node) error {
	address := node.GetAddress()
	if _, ok := s.nodes[address]; ok {
		return fmt.Errorf("node with address %v already exists in simulation", address)
	}
	s.nodes[address] = node
	s.messageQueue[address] = make(chan MessageTriplet, s.options.BufferSize)
	s.timerQueue[address] = make(chan TimerTriplet, s.options.BufferSize)
	s.interruptQueue[address] = make(chan InterruptTriplet, s.options.BufferSize)
	return nil
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
func (s *LocalSimulation) _handleMessage(ctx context.Context, node Node, mt MessageTriplet) bool {
	if node.GetState() != Running {
		return false
	}
	return node.HandleMessage(ctx, mt.Message, mt.From)
}

// handleMessage handles a message by sending it to the appropriate node.
//
// If the root node does not exist, the message is dropped.
func (s *LocalSimulation) handleMessage(ctx context.Context, mt MessageTriplet) bool {
	if _, ok := s.nodes[mt.To]; !ok {
		return false
	}
	s.LogHandleMessage(mt.From, mt.To, mt.Message)
	if s._handleMessage(ctx, s.nodes[mt.To], mt) {
		return true
	}
	for _, node := range s.nodes[mt.To].GetSubNodes() {
		if s._handleMessage(ctx, node, mt) {
			return true
		}
	}
	return false
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
func (s *LocalSimulation) _handleTimer(ctx context.Context, node Node, tt TimerTriplet) bool {
	if node.GetState() != Running {
		return false
	}
	return node.HandleTimer(ctx, tt.Timer, tt.Duration)
}

// handleTimer handles a timer by sending it to the appropriate node.
//
// If the node does not exist, the timer is dropped.
func (s *LocalSimulation) handleTimer(ctx context.Context, tt TimerTriplet) bool {
	if _, ok := s.nodes[tt.To]; !ok {
		return false
	}
	s.LogHandleTimer(tt.To, tt.Timer, tt.Duration)
	if s._handleTimer(ctx, s.nodes[tt.To], tt) {
		return true
	}
	for _, node := range s.nodes[tt.To].GetSubNodes() {
		if s._handleTimer(ctx, node, tt) {
			return true
		}
	}
	return false
}

// dropTimer drops a timer.
//
// This means the timer is not handled by any node.
func (s *LocalSimulation) dropTimer(ctx context.Context, tt TimerTriplet) {
	s.LogDropTimer(tt.To, tt.Timer, tt.Duration)
}

// _handleInterrupt handles an interrupt once the appropriate node is found.
//
// If the node is not running, the interrupt is dropped..
func (s *LocalSimulation) _handleInterrupt(ctx context.Context, node Node, it InterruptTriplet) bool {
	if node.GetState() != Running {
		return false
	}
	return node.HandleInterrupt(ctx, it.Interrupt, it.From)
}

// handleInterrupt handles an interrupt by sending it to the appropriate node.
//
// If the node does not exist, the interrupt is dropped.
func (s *LocalSimulation) handleInterrupt(ctx context.Context, it InterruptTriplet) bool {
	if _, ok := s.nodes[it.To]; !ok {
		return false
	}
	s.LogHandleInterrupt(it.From, it.To, it.Interrupt)
	if s._handleInterrupt(ctx, s.nodes[it.To], it) {
		s.LogNodeState(s.nodes[it.To])
		return true
	}
	for _, node := range s.nodes[it.To].GetSubNodes() {
		if s._handleInterrupt(ctx, node, it) {
			s.LogNodeState(node)
			return true
		}
	}
	return false
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
	for _, subNode := range node.GetSubNodes() {
		s.initNode(ctx, subNode)
	}
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
	s.startSim(ctx)
	<-ctx.Done()
	s.stopSim()
	s.generateUmlImage()
}

// generateUmlImage generates a UML image of the simulation using PlantUML (requires java).
func (s *LocalSimulation) generateUmlImage() error {
	javaPath := s.options.JavaPath
	plantumlPath := s.options.PlantumlPath
	if javaPath == "" || plantumlPath == "" {
		return fmt.Errorf("javaPath or plantumlPath not set. UML image not generated")
	}

	if s.options.UmlLogPath == "" {
		return fmt.Errorf("umlLogPath not set. UML image not generated")
	}

	cmd := exec.Command(javaPath, "-jar", plantumlPath, s.options.UmlLogPath)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
