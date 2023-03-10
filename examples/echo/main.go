package main

import (
	ds "disse/lib"
	"fmt"
	"time"
)

func main() {
	echoServerAddress, echoClientAddress := ds.Address("EchoServer"), ds.Address("EchoClient")
	pingServerAddress, pingClientAddress := echoServerAddress.SubAddress("PingServer"), echoClientAddress.SubAddress("PingClient")
	pingMessage, pongMessage, echoMessage := ds.Message("Ping"), ds.Message("Pong"), ds.Message("Echo")

	sim := ds.NewSimulation()
	pingServer := &PingServer{
		BaseNode:    ds.NewBaseNode(sim, pingServerAddress),
		pingMessage: pingMessage,
		pongMessage: pongMessage,
		PingCounter: 0,
	}
	pingClient := &PingClient{
		BaseNode:      ds.NewBaseNode(sim, pingClientAddress),
		pingMessage:   pingMessage,
		pongMessage:   pongMessage,
		serverAddress: pingServerAddress,
		pingInterval:  1 * time.Second,
		PongCounter:   0,
	}
	echoSever := &EchoServer{
		BaseNode:    ds.NewBaseNode(sim, echoServerAddress),
		echoMessage: echoMessage,
		EchoCounter: 0,
	}
	echoClient := &EchoClient{
		BaseNode:          ds.NewBaseNode(sim, echoClientAddress),
		echoInterval:      2 * time.Second,
		echoServerAddress: echoServerAddress,
		echoMessage:       echoMessage,
		EchoCounter:       0,
	}

	echoClient.AddSubNode(pingClientAddress, pingClient)
	echoSever.AddSubNode(pingServerAddress, pingServer)
	sim.AddNode(echoServerAddress, echoSever)
	sim.AddNode(echoClientAddress, echoClient)
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
	fmt.Println("Test 3: EchoClient.EchoCounter == EchoServer.EchoCounter")
	if echoClient.EchoCounter != echoSever.EchoCounter {
		fmt.Printf("Failed: EchoClient.EchoCounter (%d) != EchoServer.EchoCounter (%d)\n", echoClient.EchoCounter, echoSever.EchoCounter)
		failedTests++
	} else {
		fmt.Printf("Success: EchoClient.EchoCounter (%d) == EchoServer.EchoCounter (%d)\n", echoClient.EchoCounter, echoSever.EchoCounter)
		successTests++
	}
	fmt.Println()
	fmt.Println("Test 4: EchoClient.EchoCounter >= expected")
	expectedEchoCounter := int(sim.Duration/echoClient.echoInterval) - 1
	if echoClient.EchoCounter < expectedEchoCounter {
		fmt.Printf("Failed: EchoClient.EchoCounter (%d) < expected (%d)\n", echoClient.EchoCounter, expectedEchoCounter)
		failedTests++
	} else {
		fmt.Printf("Success: EchoClient.EchoCounter (%d) >= expected (%d)\n", echoClient.EchoCounter, expectedEchoCounter)
		successTests++
	}
	fmt.Println()
	fmt.Printf("Tests (%d): %d failed, %d success\n", failedTests+successTests, failedTests, successTests)
}
