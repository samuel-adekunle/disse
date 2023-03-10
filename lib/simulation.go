package lib

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
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

func init() {
	const (
		defaultLogFileName = ""
		logFileNameUsage   = "path to log file"
	)
	flag.StringVar(&logFileName, "logfile", defaultLogFileName, "path to log file")
	flag.StringVar(&logFileName, "l", defaultLogFileName, logFileNameUsage+" (shorthand)")
}

type Simulation struct {
	Nodes        map[Address]Node
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
		Nodes:        make(map[Address]Node),
		MessageQueue: make(chan MessageTriplet, bufferSizes.MessageBufferSize),
		TimerQueue:   make(chan TimerTriplet, bufferSizes.TimerBufferSize),
		MinLatency:   defaultMinLatency,
		MaxLatency:   defaultMaxLatency,
		Duration:     defaultDuration,
	}
}

func (s *Simulation) AddNode(address Address, node Node) {
	s.Nodes[address] = node
}

func (s *Simulation) RandomLatency() time.Duration {
	return s.MinLatency + time.Duration(rand.Int63n(int64(s.MaxLatency-s.MinLatency)))
}

func (s *Simulation) Run() {
	flag.Parse()
	var logFile *os.File = os.Stdout
	if logFileName != "" {
		var err error
		logFile, err = os.Create(logFileName)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()
	}

	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
	log.SetOutput(logFile)
	log.Printf("SetLogOutput(%v)\n", logFile.Name())
	log.Printf("StartSim(%v)\n", s.Duration)

	var ctx context.Context
	if s.Duration == Infinity {
		ctx = context.Background()
	} else {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), s.Duration)
		defer cancel()
	}

	for _, node := range s.Nodes {
		node.Init(ctx)
	}

	var wg sync.WaitGroup
	for {
		select {
		case <-ctx.Done():
			if s.Duration != Infinity {
				log.Println("StopSim()")
				close(s.MessageQueue)
				close(s.TimerQueue)
				wg.Wait()
				log.Println("EndSim()")
				return
			}
		case mt := <-s.MessageQueue:
			wg.Add(1)
			go func() {
				time.Sleep(s.RandomLatency())
				s.Nodes[mt.To.Root()].HandleMessage(ctx, mt.Message, mt.From)
				wg.Done()
			}()
		case tt := <-s.TimerQueue:
			wg.Add(1)
			go func() {
				time.Sleep(tt.Length)
				s.Nodes[tt.From.Root()].HandleTimer(ctx, tt.Timer, tt.Length)
				wg.Done()
			}()
		}
	}
}
