package main

import (
	"testing"
)

func TestBeb(t *testing.T) {
	sim.Run()
	t.Run("TestSentEqualsReceived", func(t *testing.T) {
		for _, client := range bebClients {
			if len(client.(*BebClient).Received) != len(bebServer.Sent) {
				t.Errorf("Client %s received %d messages, but server sent %d messages", client.GetAddress(), len(client.(*BebClient).Received), len(bebServer.Sent))
			}
		}
	})
}
