package disse

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

const (
	Infinity                 = time.Duration(0)
	DefaultMessageBufferSize = 100
	DefaultTimerBufferSize   = 100
	DefaultMinLatency        = 10 * time.Millisecond
	DefaultMaxLatency        = 100 * time.Millisecond
	DefaultDuration          = Infinity
)

const (
	javaEnvVarName     = "DISSE_JAVA_PATH"
	plantumlEnvVarName = "DISSE_PLANTUML_JAR"
)

// debugLog is the log used for debug messages.
var debugLog *log.Logger
var debugLogPath string

// umlLog is the log used for UML messages.
var umlLog *log.Logger
var umlLogPath string

// Init sets up the command line flags for the simulation executable.
// The log file name is the file where the simulation logs will be written.
// The default log file name is the name of the executable with the .log extension.
// The UML file name is the file where the UML diagram will be written.
// The default UML file name is the name of the executable with the .uml extension.
func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	workDir := strings.Split(wd, "/")
	dirName := workDir[len(workDir)-1]

	defaultLogPath := fmt.Sprintf("%s.log", dirName)
	logFileNameUsage := "path to log file"
	defaultUmlPath := fmt.Sprintf("%s.uml", dirName)
	umlFileNameUsage := "path to UML diagram file"

	flag.StringVar(&debugLogPath, "logfile", defaultLogPath, "path to log file")
	flag.StringVar(&debugLogPath, "l", defaultLogPath, logFileNameUsage+" (shorthand)")
	flag.StringVar(&umlLogPath, "umlfile", defaultUmlPath, "path to UML diagram file")
	flag.StringVar(&umlLogPath, "u", defaultUmlPath, umlFileNameUsage+" (shorthand)")
}

// Simulation sets up and runs the distributed system simulation.
type Simulation struct {
	nodes        map[Address]Node
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
//
// If the message is successfully handled, true is returned, otherwise false.
func (s *Simulation) HandleMessage(ctx context.Context, mt MessageTriplet) (handled bool) {
	if node, ok := s.nodes[mt.To]; ok {
		if node.GetState() != Running {
			return false
		}
		debugLog.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message)
		return node.HandleMessage(ctx, mt.Message, mt.From)
	}
	return s.nodes[mt.To.Root()].SubNodesHandleMessage(ctx, mt)
}

// DropMessage drops a message.
//
// This means the message is not handled by any node.
func (s *Simulation) DropMessage(ctx context.Context, mt MessageTriplet) {
	debugLog.Printf("DropMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message)
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
		debugLog.Printf("HandleTimer(-> %v, %v, %v)\n", tt.To, tt.Timer, tt.Duration)
		return node.HandleTimer(ctx, tt.Timer, tt.Duration)
	}
	return s.nodes[tt.To.Root()].SubNodesHandleTimer(ctx, tt)
}

// DropTimer drops a timer.
//
// This means the timer is not handled by any node.
func (s *Simulation) DropTimer(ctx context.Context, tt TimerTriplet) {
	debugLog.Printf("DropTimer(-> %v, %v, %v)\n", tt.To, tt.Timer, tt.Duration)
}

// HandleInterrupt handles an interrupt by sending it to the appropriate node.
//
// If the node is not running, the interrupt is dropped.
//
// If the address does not match the root node, the sub nodes are checked recursively for a match.
//
// If an unknown interrupt is received, the interrupt is dropped and the function returns false, otherwise true.
func (s *Simulation) HandleInterrupt(ctx context.Context, ip InterruptPair) bool {
	if node, ok := s.nodes[ip.To]; ok {
		if node.GetState() == Stopped {
			return false
		}
		debugLog.Printf("HandleInterrupt(-> %v, %v)\n", ip.To, ip.Interrupt)
		return node.HandleInterrupt(ctx, ip.Interrupt)
	}
	return s.nodes[ip.To.Root()].SubNodesHandleInterrupt(ctx, ip)
}

// DropInterrupt drops an interrupt.
//
// This means the interrupt is not handled by any node.
func (s *Simulation) DropInterrupt(ctx context.Context, ip InterruptPair) {
	debugLog.Printf("DropInterrupt(-> %v, %v)\n", ip.To, ip.Interrupt)
}

// startSim starts the simulation by initializing all nodes and sub nodes.
func (s *Simulation) startSim(ctx context.Context) {
	debugLog.Printf("StartSim(%v)\n", s.Duration)
	for address, node := range s.nodes {
		debugLog.Printf("Init(%v)\n", address)
		node.Init(ctx)
		node.SubNodesInit(ctx)
	}
}

// stopSim stops the simulation by closing the message and timer queues and waiting for all nodes to stop doing work.
func (s *Simulation) stopSim() {
	debugLog.Println("StopSim()")
	close(s.messageQueue)
	close(s.timerQueue)
}

// generateUmlImage generates a UML image of the simulation using PlantUML (requires java).
func (s *Simulation) generateUmlImage() {
	javaPath := os.Getenv(javaEnvVarName)
	plantumlPath := os.Getenv(plantumlEnvVarName)

	if javaPath == "" || plantumlPath == "" {
		debugLog.Printf("javaPath (%v) or plantumlPath (%v) not set. UML image not generated.\n", javaPath, plantumlPath)
		return
	}

	if umlLogPath == os.DevNull {
		debugLog.Printf("umlPath (%v) set to /dev/null. UML image not generated.\n", umlLogPath)
		return
	}

	cmd := exec.Command(javaPath, "-jar", plantumlPath, umlLogPath)
	err := cmd.Run()
	if err != nil {
		debugLog.Printf("Error %v when generating UML image. Check if javaPath (%v) and plantumlPath (%v) are correctly set.\n", err, javaPath, plantumlPath)
	}
}

// Run runs the simulation.
//
// The simulation run by polling the message and timer queues and sending the messages and timers to the appropriate nodes.
func (s *Simulation) Run() {
	flag.Parse()
	godotenv.Load()
	logFile, err := os.Create(debugLogPath)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	debugLog = log.New(logFile, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	umlFile, err := os.Create(umlLogPath)
	if err != nil {
		panic(err)
	}
	defer s.generateUmlImage()
	defer umlFile.Close()
	umlLog = log.New(umlFile, "", 0)
	umlLog.Println("@startuml")
	umlLog.Println("!theme reddress-lightred")
	umlLog.Println("skinparam shadowing false")
	umlLog.Println("skinparam sequenceArrowThickness 1")
	umlLog.Println("skinparam responseMessageBelowArrow true")
	umlLog.Println("skinparam sequenceMessageAlign right")
	defer umlLog.Println("@enduml")

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
				debugLog.Println("EndSim()")
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
