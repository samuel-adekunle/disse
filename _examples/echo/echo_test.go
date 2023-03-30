package echo

import (
	ds "disse"
	"testing"
	"time"
)

var sim *ds.Simulation
var pingServer *PingServer
var pingClient *PingClient
var echoServer *EchoServer
var echoClient *EchoClient

func newSim() {
	echoServerAddress, echoClientAddress := ds.Address("EchoServer"), ds.Address("EchoClient")
	pingServerAddress, pingClientAddress := echoServerAddress.SubAddress("PingServer"), echoClientAddress.SubAddress("PingClient")
	pingMessage, pongMessage := ds.NewMessage(ds.MessageId("Ping"), nil), ds.NewMessage(ds.MessageId("Pong"), nil)
	echoMessage := ds.NewMessage(ds.MessageId("Echo"), nil)

	sim = ds.NewSimulationWithBuffer(&ds.BufferSizes{
		MessageBufferSize: 10,
		TimerBufferSize:   10,
	})
	pingServer = &PingServer{
		BaseNode:    ds.NewBaseNode(sim, pingServerAddress),
		pingMessage: pingMessage,
		pongMessage: pongMessage,
		PingCounter: 0,
	}
	pingClient = &PingClient{
		BaseNode:      ds.NewBaseNode(sim, pingClientAddress),
		pingMessage:   pingMessage,
		pongMessage:   pongMessage,
		serverAddress: pingServerAddress,
		pingInterval:  200 * time.Millisecond,
		PongCounter:   0,
	}
	echoServer = &EchoServer{
		BaseNode:    ds.NewBaseNode(sim, echoServerAddress),
		echoMessage: echoMessage,
		EchoCounter: 0,
	}
	echoClient = &EchoClient{
		BaseNode:          ds.NewBaseNode(sim, echoClientAddress),
		echoInterval:      500 * time.Millisecond,
		echoServerAddress: echoServerAddress,
		echoMessage:       echoMessage,
		EchoCounter:       0,
	}

	echoClient.AddSubNode(pingClientAddress, pingClient)
	echoServer.AddSubNode(pingServerAddress, pingServer)
	sim.AddNode(echoServerAddress, echoServer)
	sim.AddNode(echoClientAddress, echoClient)
	sim.Duration = 2 * time.Second
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
	pingClient.pingInterval = 20 * time.Millisecond
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

func TestExpectedEchoCounter(t *testing.T) {
	newSim()
	sim.Run()

	expectedEchoCounter := int(sim.Duration/echoClient.echoInterval) - 1
	if echoClient.EchoCounter < expectedEchoCounter {
		t.Errorf("EchoClient.EchoCounter (%d) less than expected value of (%d)", echoClient.EchoCounter, expectedEchoCounter)
	}
}

func TestEchoCounterEqualEchoCounter(t *testing.T) {
	newSim()
	sim.Run()

	if echoClient.EchoCounter != echoServer.EchoCounter {
		t.Errorf("EchoClient.EchoCounter (%d) != EchoServer.EchoCounter (%d)", echoClient.EchoCounter, echoServer.EchoCounter)
	}
}

func TestEchoFlooding(t *testing.T) {
	newSim()
	echoClient.echoInterval = 50 * time.Millisecond
	sim.Run()

	expectedNumberOfEchos := int(sim.Duration/echoClient.echoInterval) - 1
	expectedDroppedEchos, expectedEchoLag := 3, 1
	if echoClient.EchoCounter < echoServer.EchoCounter-expectedEchoLag {
		t.Errorf("EchoClient.EchoCounter (%d) less than EchoServer.EchoCounter (%d) by more than the expected echo lag (%v)", echoClient.EchoCounter, echoServer.EchoCounter, expectedEchoLag)
	}
	if echoServer.EchoCounter < expectedNumberOfEchos-expectedDroppedEchos {
		t.Errorf("EchoServer.EchoCounter (%d) less than expected number of echos (%d) by more than expected number of dropped echos (%v)", echoServer.EchoCounter, expectedNumberOfEchos, expectedDroppedEchos)
	}
}
