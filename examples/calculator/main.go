package main

import (
	ds "github.com/samuel-adekunle/disse"
)

func main() {
	sim := ds.NewLocalSimulation(nil)

	adderAddress := ds.Address("adder")
	adderNode := &AdderNode{
		LocalNode: ds.NewLocalNode(sim, adderAddress),
	}

	multiplierAddress := ds.Address("multiplier")
	multiplierNode := &MultiplierNode{
		LocalNode: ds.NewLocalNode(sim, multiplierAddress),
	}

	calculatorAddress := ds.Address("calculator")
	calculatorNode := &CalculatorNode{
		LocalNode: ds.NewLocalNode(sim, calculatorAddress),
	}
	calculatorNode.AddSubNode(adderNode)
	calculatorNode.AddSubNode(multiplierNode)
	sim.AddNode(calculatorNode)

	testAddress := ds.Address("test")
	testNode := &TestNode{
		LocalNode:  ds.NewLocalNode(sim, testAddress),
		A:          10,
		B:          5,
		calculator: calculatorAddress,
	}
	sim.AddNode(testNode)

	sim.Run()
}
