package lib

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"
)

var logFileName string

func init() {
	flag.StringVar(&logFileName, "log", "", "path to log file")
}

type Simulation struct {
	Nodes        map[Address]Node
	MessageQueue chan MessageTriplet
	TimerQueue   chan TimerTriplet
	MinLatency   time.Duration
	MaxLatency   time.Duration
	Duration     time.Duration
}

const Infinity time.Duration = 0

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

	var wg sync.WaitGroup

	log.Printf("StartSim(%v)\n", s.Duration)

	for _, node := range s.Nodes {
		node.Init()
	}

	for {
		select {
		case mt := <-s.MessageQueue:
			wg.Add(1)
			go func() {
				time.Sleep(s.RandomLatency())
				s.Nodes[mt.To].HandleMessage(mt.Message, mt.From)
				wg.Done()
			}()
		case tt := <-s.TimerQueue:
			wg.Add(1)
			go func() {
				time.Sleep(tt.Length)
				s.Nodes[tt.From].HandleTimer(tt.Timer, tt.Length)
				wg.Done()
			}()
		case <-time.After(s.Duration):
			if s.Duration != Infinity {
				wg.Wait()
				log.Println("StopSim()")
				return
			}
		}
	}
}
