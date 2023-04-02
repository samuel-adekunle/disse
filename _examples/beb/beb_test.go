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
func TestFaultyNode(t *testing.T) {
	initBebSimulation()
	helloNode := helloNodes[0].(*HelloNode)
	faultyHelloNode := &FaultyHelloNode{
		HelloNode:  helloNode,
		faultAfter: 1500 * time.Millisecond,
	}
	sim.AddNode(helloNode.GetAddress(), faultyHelloNode)
	sim.Run()
}
