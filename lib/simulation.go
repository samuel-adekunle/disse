package lib

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
	defaultMessageBufferSize = 100
	defaultTimerBufferSize   = 100
	defaultFaultBufferSize   = 100
	defaultMinLatency        = 10 * time.Millisecond
	defaultMaxLatency        = 100 * time.Millisecond
	defaultDuration          = Infinity
)

var logFileName string
var umlFileName string

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

type Simulation struct {
	nodes        map[Address]Node
	debugLog     *log.Logger
	umlLog       *log.Logger
	MessageQueue chan MessageTriplet
	TimerQueue   chan TimerTriplet
	Faults       []FaultTriplet
	MinLatency   time.Duration
	MaxLatency   time.Duration
	Duration     time.Duration
}

type BufferSizes struct {
	MessageBufferSize int
	TimerBufferSize   int
	FaultBufferSize   int
}

func NewSimulation() *Simulation {
	return NewSimulationWithBuffer(nil)
}

func NewSimulationWithBuffer(bufferSizes *BufferSizes) *Simulation {
	if bufferSizes == nil {
		bufferSizes = &BufferSizes{
			MessageBufferSize: defaultMessageBufferSize,
			TimerBufferSize:   defaultTimerBufferSize,
			FaultBufferSize:   defaultFaultBufferSize,
		}
	}
	return &Simulation{
		nodes:        make(map[Address]Node),
		MessageQueue: make(chan MessageTriplet, bufferSizes.MessageBufferSize),
		TimerQueue:   make(chan TimerTriplet, bufferSizes.TimerBufferSize),
		Faults:       make([]FaultTriplet, 0, bufferSizes.FaultBufferSize),
		MinLatency:   defaultMinLatency,
		MaxLatency:   defaultMaxLatency,
		Duration:     defaultDuration,
	}
}

func (s *Simulation) AddNode(address Address, node Node) {
	s.nodes[address] = node
}

func (s *Simulation) AddNodes(addresses []Address, nodes []Node) {
	if len(addresses) != len(nodes) {
		panic("addresses and nodes must have the same length")
	}
	for i := range addresses {
		s.AddNode(addresses[i], nodes[i])
	}
}

func (s *Simulation) AddFault(fault FaultTriplet) {
	s.Faults = append(s.Faults, fault)
}

func (s *Simulation) AddFaults(faults []FaultTriplet) {
	s.Faults = append(s.Faults, faults...)
}

func (s *Simulation) randomLatency() time.Duration {
	return s.MinLatency + time.Duration(rand.Int63n(int64(s.MaxLatency-s.MinLatency)))
}

func (s *Simulation) HandleMessage(ctx context.Context, mt MessageTriplet) {
	if node, ok := s.nodes[mt.To]; ok {
		s.debugLog.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message)
		node.HandleMessage(ctx, mt.Message, mt.From)
	} else {
		s.nodes[mt.To.Root()].SubNodesHandleMessage(ctx, mt)
	}
}

func (s *Simulation) HandleTimer(ctx context.Context, tt TimerTriplet) {
	if node, ok := s.nodes[tt.To]; ok {
		s.debugLog.Printf("HandleTimer(%v, %v, %v)\n", tt.To, tt.Timer, tt.Duration)
		node.HandleTimer(ctx, tt.Timer, tt.Duration)
	} else {
		s.nodes[tt.To.Root()].SubNodesHandleTimer(ctx, tt)
	}
}

func (s *Simulation) HandleFault(ctx context.Context, fault FaultTriplet) {
	switch fault.Fault.Name {
	case Stop:
		// TODO: Stop all subnodes
	case Resume:
		// TODO: Resume all subnodes
	case Restart:
		// TODO: Restart all subnodes
	case Sleep:
		// TODO: Sleep all subnodes
	}
}

func (s *Simulation) startSim(ctx context.Context) {
	s.debugLog.Printf("StartSim(%v)\n", s.Duration)
	for address, node := range s.nodes {
		s.debugLog.Printf("Init(%v)\n", address)
		node.Init(ctx)
		node.SubNodesInit(ctx)
	}
}

func (s *Simulation) stopSim() {
	s.debugLog.Println("StopSim()")
	close(s.MessageQueue)
	close(s.TimerQueue)
}

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

	for _, ft := range s.Faults {
		wg.Add(1)
		go func(_ft FaultTriplet) {
			time.Sleep(_ft.After)
			s.HandleFault(ctx, _ft)
			wg.Done()
		}(ft)
	}

	for {
		select {
		case <-ctx.Done():
			if s.Duration != Infinity {
				s.stopSim()
				wg.Wait()
				s.debugLog.Println("EndSim()")
				return
			}
		case mt := <-s.MessageQueue:
			wg.Add(1)
			go func() {
				time.Sleep(s.randomLatency())
				s.HandleMessage(ctx, mt)
				wg.Done()
			}()
		case tt := <-s.TimerQueue:
			wg.Add(1)
			go func() {
				time.Sleep(tt.Duration)
				s.HandleTimer(ctx, tt)
				wg.Done()
			}()
		}
	}
}
