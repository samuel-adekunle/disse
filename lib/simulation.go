package lib

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Simulation struct {
	Nodes        map[Address]Node
	MessageQueue chan MessageTriplet
	TimerQueue   chan TimerTriplet
	MinLatency   time.Duration
	MaxLatency   time.Duration
}

const Infinity time.Duration = 0

func (s *Simulation) RandomLatency() time.Duration {
	return s.MinLatency + time.Duration(rand.Int63n(int64(s.MaxLatency-s.MinLatency)))
}

func (s *Simulation) Run(duration time.Duration) {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)

	wg := &sync.WaitGroup{}

	for _, node := range s.Nodes {
		node.Init()
	}

	for {
		select {
		case mt := <-s.MessageQueue:
			wg.Add(1)
			go func() {
				time.Sleep(s.RandomLatency())
				log.Printf("HandleMessage(%v -> %v, %v)\n", mt.From, mt.To, mt.Message)
				s.Nodes[mt.To].HandleMessage(mt.Message, mt.From)
				wg.Done()
			}()
		case tt := <-s.TimerQueue:
			wg.Add(1)
			go func() {
				time.Sleep(tt.Length)
				log.Printf("HandleTimer(%v, %v, %v)\n", tt.From, tt.Timer, tt.Length)
				s.Nodes[tt.From].HandleTimer(tt.Timer)
				wg.Done()
			}()
		case <-time.After(duration):
			if duration != Infinity {
				wg.Wait()
				log.Println("Simulation finished")
				return
			}
		}
	}
}
