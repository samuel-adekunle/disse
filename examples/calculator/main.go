package main

import (
	ds "github.com/samuel-adekunle/disse"
)

func main() {
	sim := ds.NewLocalSimulation(nil)

	adderAddress := ds.Address("calculator.adder")
	adderNode := &AdderNode{
		LocalNode: ds.NewLocalNode(sim, adderAddress),
	}

	multiplierAddress := ds.Address("calculator.adder.multiplier")
	multiplierNode := &MultiplierNode{
		LocalNode: ds.NewLocalNode(sim, multiplierAddress),
	}

	calculatorAddress := ds.Address("calculator")
	calculatorNode := &CalculatorNode{
		LocalNode: ds.NewLocalNode(sim, calculatorAddress),
	}
	sim.AddNode(calculatorNode)
	calculatorNode.AddSubNode(adderNode)
	adderNode.AddSubNode(multiplierNode)

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
