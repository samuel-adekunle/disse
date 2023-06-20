package main

import (
	"fmt"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

const SIM_TIME = 5 * time.Second
const REPEAT = 5

func main() {
	// uncomment to profile memory usage
	// defer profile.Start(profile.ProfilePath("."), profile.MemProfile).Stop()
	var numNodes = []int{50, 100, 200, 500, 1000, 1500, 2000, 3000, 5000, 7500, 10000}
	var simDelta time.Duration
	fmt.Println("numNodes\tsimDelta")
	for _, NUM_NODES := range numNodes {
		for i := 0; i < REPEAT; i++ {
			sim := ds.NewLocalSimulation(&ds.LocalSimulationOptions{
				MinLatency:   10 * time.Millisecond,
				MaxLatency:   100 * time.Millisecond,
				Duration:     SIM_TIME,
				BufferSize:   NUM_NODES,
				DebugLogPath: "debug.log",
				UmlLogPath:   "uml.log",
				JavaPath:     ds.DefaultJavaPath,
				PlantumlPath: ds.DefaultPlantumlPath,
			})

			addresses := make([]ds.Address, NUM_NODES)
			for i := 0; i < NUM_NODES; i++ {
				addresses[i] = ds.Address(fmt.Sprintf("pingNode%d", i))
			}

			for i := 0; i < NUM_NODES; i++ {
				node := &PingNode{
					LocalNode: ds.NewLocalNode(sim, addresses[i]),
					Nodes:     addresses,
				}
				sim.AddNode(node)
			}

			start := time.Now()
			sim.Run()
			realSimTime := time.Since(start)
			simDelta += realSimTime - SIM_TIME
		}
		fmt.Printf("%v\t%v\n", NUM_NODES, simDelta/REPEAT)
	}
}
