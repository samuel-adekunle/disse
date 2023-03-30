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
	MinLatency   time.Duration
	MaxLatency   time.Duration
	Duration     time.Duration
}

type BufferSizes struct {
	MessageBufferSize int
	TimerBufferSize   int
}

func NewSimulation() *Simulation {
	return NewSimulationWithBuffer(nil)
}

func NewSimulationWithBuffer(bufferSizes *BufferSizes) *Simulation {
	if bufferSizes == nil {
		bufferSizes = &BufferSizes{
			MessageBufferSize: defaultMessageBufferSize,
			TimerBufferSize:   defaultTimerBufferSize,
		}
	}
	return &Simulation{
		nodes:        make(map[Address]Node),
		MessageQueue: make(chan MessageTriplet, bufferSizes.MessageBufferSize),
		TimerQueue:   make(chan TimerTriplet, bufferSizes.TimerBufferSize),
		MinLatency:   defaultMinLatency,
		MaxLatency:   defaultMaxLatency,
		Duration:     defaultDuration,
	}
}

func (s *Simulation) randomLatency() time.Duration {
	return s.MinLatency + time.Duration(rand.Int63n(int64(s.MaxLatency-s.MinLatency)))
}

func (s *Simulation) AddNode(address Address, node Node) {
	s.nodes[address] = node
}

func (s *Simulation) AddNodes(addresses []Address, nodes []Node) (err error) {
	if len(addresses) != len(nodes) {
		return fmt.Errorf("length of addresses (%v) does not match length of nodes (%v)", len(addresses), len(nodes))
	}
	for i := range addresses {
		s.AddNode(addresses[i], nodes[i])
	}
	return nil
}

func (s *Simulation) HandleMessage(ctx context.Context, mt MessageTriplet) {
	if node, ok := s.nodes[mt.To]; ok {
		if node.GetState() != Running {
			s.debugLog.Printf("DroppedMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message.Id)
			return
		}
		s.debugLog.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message.Id)
		node.HandleMessage(ctx, mt.Message, mt.From)
	} else {
		s.nodes[mt.To.Root()].SubNodesHandleMessage(ctx, mt)
	}
}

func (s *Simulation) HandleTimer(ctx context.Context, tt TimerTriplet) {
	if node, ok := s.nodes[tt.To]; ok {
		if node.GetState() != Running {
			s.debugLog.Printf("DroppedTimer(%v, %v, %v)\n", tt.To, tt.Timer.Id, tt.Duration)
			return
		}
		s.debugLog.Printf("HandleTimer(%v, %v, %v)\n", tt.To, tt.Timer.Id, tt.Duration)
		node.HandleTimer(ctx, tt.Timer, tt.Duration)
	} else {
		s.nodes[tt.To.Root()].SubNodesHandleTimer(ctx, tt)
	}
}

func (s *Simulation) HandleInterrupt(ctx context.Context, ip InterruptPair) (err error) {
	if node, ok := s.nodes[ip.To]; ok {
		if node.GetState() == Stopped {
			s.debugLog.Printf("DroppedInterrupt(%v, %v)\n", ip.To, ip.Interrupt.Id)
			return
		}
		s.debugLog.Printf("HandleInterrupt(%v, %v)\n", ip.To, ip.Interrupt.Id)
		err = node.HandleInterrupt(ctx, ip.Interrupt)
	} else {
		err = s.nodes[ip.To.Root()].SubNodesHandleInterrupt(ctx, ip)
	}
	return
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
