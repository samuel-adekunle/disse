package main

import (
	ds "github.com/samuel-adekunle/disse"
)

func main() {
	sim := ds.NewLocalSimulation(nil)

	calculatorAddress := ds.Address("calculator")
	adderAddress := calculatorAddress.NewSubAddress("adder")
	multiplierAddress := adderAddress.NewSubAddress("multiplier")

	adderNode := &AdderNode{
		LocalNode: ds.NewLocalNode(sim, adderAddress),
	}
	multiplierNode := &MultiplierNode{
		LocalNode: ds.NewLocalNode(sim, multiplierAddress),
	}
	calculatorNode := &CalculatorNode{
		LocalNode: ds.NewLocalNode(sim, calculatorAddress),
	}

	sim.AddNode(calculatorNode)
	adderNode.AddSubNode(multiplierNode)
	calculatorNode.AddSubNode(adderNode)

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
