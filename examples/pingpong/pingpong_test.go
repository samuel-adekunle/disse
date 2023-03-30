package pingpong

import (
	ds "disse/lib"
	"testing"
	"time"
)

var sim *ds.Simulation
var pingServer *PingServer
var pingClient *PingClient

func newSim() {
	pingMessage, pongMessage := ds.NewMessage(ds.MessageId("Ping"), nil), ds.NewMessage(ds.MessageId("Pong"), nil)
	serverAddress, clientAddress := ds.Address("PingServer"), ds.Address("PingClient")
	sim = ds.NewSimulationWithBuffer(&ds.BufferSizes{
		MessageBufferSize: 5,
		TimerBufferSize:   5,
	})
	pingServer = &PingServer{
		BaseNode:    ds.NewBaseNode(sim, serverAddress),
		pingMessage: pingMessage,
		pongMessage: pongMessage,
		PingCounter: 0,
	}
	pingClient = &PingClient{
		BaseNode:      ds.NewBaseNode(sim, clientAddress),
		pingMessage:   pingMessage,
		pongMessage:   pongMessage,
		serverAddress: serverAddress,
		pingInterval:  200 * time.Millisecond,
		PongCounter:   0,
	}
	sim.AddNode(serverAddress, pingServer)
	sim.AddNode(clientAddress, pingClient)
	sim.Duration = 1 * time.Second
	sim.MinLatency = 10 * time.Millisecond
	sim.MaxLatency = 100 * time.Millisecond
}

func TestExpectedPingCounter(t *testing.T) {
	newSim()
	sim.Run()

	expectedPongCounter := int(sim.Duration/pingClient.pingInterval) - 1
	if pingClient.PongCounter < expectedPongCounter {
		t.Errorf("PingClient.PongCounter (%d) less than expected value of (%d)", pingClient.PongCounter, expectedPongCounter)
	}
}

func TestPingCounterEqualPongCounter(t *testing.T) {
	newSim()
	sim.Run()

	if pingClient.PongCounter != pingServer.PingCounter {
		t.Errorf("PingClient.PongCounter (%d) != PingServer.PingCounter (%d)", pingClient.PongCounter, pingServer.PingCounter)
	}
}

func TestPingFlooding(t *testing.T) {
	newSim()
	pingClient.pingInterval = 10 * time.Millisecond
	sim.Run()

	expectedNumberOfPings := int(sim.Duration/pingClient.pingInterval) - 1
	expectedDroppedPings, expectedPongLag := 10, 5
	if pingClient.PongCounter < pingServer.PingCounter-5 {
		t.Errorf("PingClient.PongCounter (%d) less than PingServer.PingCounter (%d) by more than the expected pong lag (%v)", pingClient.PongCounter, pingServer.PingCounter, expectedPongLag)
	}
	if pingServer.PingCounter < expectedNumberOfPings-expectedDroppedPings {
		t.Errorf("PingServer.PingCounter (%d) less than expected number of pings (%d) by more than expected number of dropped pings (%v)", pingServer.PingCounter, expectedNumberOfPings, expectedDroppedPings)
	}
}
