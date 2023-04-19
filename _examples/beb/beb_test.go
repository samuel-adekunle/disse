package main

import (
	"testing"
	"time"

	ds "github.com/samuel-adekunle/disse"
)

// TestBeb tests the broadcast reliable broadcast module under normal conditions.
//
// The test is successful if the following conditions are met:
//  1. All messages are handled by a correct node.
//  2. No message is handled more than once by a correct node.
//  3. All messages are created by a correct node.
func TestBeb(t *testing.T) {
	initBebSimulation()
	sim.Run()
	t.Run("TestValidity", testValidity)
	t.Run("TestNoDuplication", testNoDuplication)
	t.Run("TestNoCreation", testNoCreation)
}

// testValidity tests that all messages are handled by a correct node.
func testValidity(t *testing.T) {
	checkValidity := func(message ds.MessageId, handledMessages map[ds.MessageId]int) {
		if _, ok := handledMessages[message]; !ok {
			t.Errorf("Message was not handled by a correct node")
		}
	}

	for _, helloNode := range helloNodes {
		helloNode := helloNode.(*HelloNode)
		for _, message := range helloNode.sentMessages {
			checkValidity(message, beb.handledMessages)

			for _, node := range helloNodes {
				node := node.(*HelloNode)
				checkValidity(message, node.handledMessages)
			}
		}
	}
}

// testNoDuplication tests that no message is handled more than once by a correct node.
func testNoDuplication(t *testing.T) {
	checkDuplicates := func(handledMessages map[ds.MessageId]int) {
		for _, count := range handledMessages {
			if count > 1 {
				t.Errorf("Message was handled more than once by a correct node")
			}
		}
	}

	checkDuplicates(beb.handledMessages)

	for _, helloNode := range helloNodes {
		helloNode := helloNode.(*HelloNode)
		checkDuplicates(helloNode.handledMessages)
	}
}

// testNoCreation tests that all messages are created by a correct node.
func testNoCreation(t *testing.T) {
	checkCreated := func(message ds.MessageId) {
		created := false
		for _, helloNode := range helloNodes {
			helloNode := helloNode.(*HelloNode)
			if _, ok := helloNode.handledMessages[message]; ok {
				created = true
				break
			}
		}
		if !created {
			t.Errorf("Message was not created by a correct node")
		}
	}
	for message := range beb.handledMessages {
		checkCreated(message)
	}
	for _, helloNode := range helloNodes {
		helloNode := helloNode.(*HelloNode)
		for message := range helloNode.handledMessages {
			checkCreated(message)
		}
	}
}

// TestFaultyNode tests the broadcast reliable broadcast module with a faulty node.
//
// The test is successful if the following conditions are met:
//  1. All messages are handled by a correct node.
//  2. No message is handled more than once by a correct node.
//  3. All messages are created by a correct node.
//
// The faulty node stops responding after 1500 milliseconds and drops all messages after that.
//
// We expect that the faulty node will not be able to handle / create any more messages after 1500 milliseconds.
func TestFaultyNode(t *testing.T) {
	initBebSimulation()
	faultyNode := helloNodes[0].(*HelloNode)
	faultyNodeAddress := faultyNode.GetAddress()
	faultyHelloNode = &FaultyHelloNode{
		HelloNode:  faultyNode,
		faultAfter: 1500 * time.Millisecond,
	}
	sim.RemoveNode(faultyNodeAddress)
	sim.AddNode(faultyNodeAddress, faultyHelloNode)
	sim.Run()
	// TODO: Add a test to confirm that the faulty node has stopped responding after 1500 milliseconds.
	// NOTE: For now, confirm via stdout / logs that the faulty node has stopped responding after 1500 milliseconds.
}
