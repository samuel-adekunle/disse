# Calculator Example

This example demonstrates how to use node composition to reuse handlers and design nodes in a modular way.

The calculator node is composed of an adder node and a multiplier node. Both are subnodes and cannot be accessed directly by other nodes.

When handling a message, the simulation checks if the root node can handle the message. If it cannot, it iterates through the subnodes and checks if they can handle the message. The first subnode that can handle the message is used to handle the message and the simulation stops searching for a handler.

In this case, the calculator root node cannot handle any message and any Calculator operation message is handled by the adder node or the multiplier node.

## Implementation
 - View the implementation [here](./calculator.go).
 - View the implementation of the adder node [here](./adder.go).
 - View the implementation of the multiplier node [here](./multiplier.go).
 - View the simulation [here](./main.go)