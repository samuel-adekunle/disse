package main

import (
	ds "disse/lib"
	"fmt"
	"time"
)

func main() {
	pingMessage, pongMessage := ds.Message("Ping"), ds.Message("Pong")
	serverAddress, clientAddress := ds.Address("PingServer"), ds.Address("PingClient")

	sim := ds.NewSimulation()
	pingServer := &PingServer{
		BaseNode:    ds.NewBaseNode(sim, serverAddress),
		pingMessage: pingMessage,
		pongMessage: pongMessage,
		PingCounter: 0,
	}
	pingClient := &PingClient{
		BaseNode:      ds.NewBaseNode(sim, clientAddress),
		pingMessage:   pingMessage,
		pongMessage:   pongMessage,
		serverAddress: serverAddress,
		pingInterval:  1 * time.Second,
		PongCounter:   0,
	}

	sim.AddNode(serverAddress, pingServer)
	sim.AddNode(clientAddress, pingClient)
	sim.Duration = 10 * time.Second
	sim.Run()

	// Testing
	fmt.Printf("Run Tests\n")
	failedTests, successTests := 0, 0

	fmt.Println("Test 1: PingClient.PongCounter == PingServer.PingCounter")
	if pingClient.PongCounter != pingServer.PingCounter {
		fmt.Printf("Failed: PingClient.PongCounter (%d) != PingServer.PingCounter (%d)\n", pingClient.PongCounter, pingServer.PingCounter)
		failedTests++
	} else {
		fmt.Printf("Success: PingClient.PongCounter (%d) == PingServer.PingCounter (%d)\n", pingClient.PongCounter, pingServer.PingCounter)
		successTests++
	}
	fmt.Println()
	fmt.Println("Test 2: PingClient.PongCounter >= expected")
	expectedPongCounter := int(sim.Duration/pingClient.pingInterval) - 1
	if pingClient.PongCounter < expectedPongCounter {
		fmt.Printf("Failed: PingClient.PongCounter (%d) < expected (%d)\n", pingClient.PongCounter, expectedPongCounter)
		failedTests++
	} else {
		fmt.Printf("Success: PingClient.PongCounter (%d) >= expected (%d)\n", pingClient.PongCounter, expectedPongCounter)
		successTests++
	}
	fmt.Println()
	fmt.Printf("Tests (%d): %d failed, %d success\n", failedTests+successTests, failedTests, successTests)
}
