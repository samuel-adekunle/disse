# Faulty Example

This example demonstrates how to use the interrupt system to simulate faulty behavior, and how to reuse the standard library components in custom modules.

The faulty node has a specified lifetime, after which it sends a `StopInterrupt` to itself which causes it to crash and never recover.

Leader election and perfect failure detector nodes also exist in the simulation to detect crashed nodes and elect a new leader if necessary.

The leader election node and perfect failure detector node are both standard library components, and are reused in this example.

By design, the faulty node is the leader of the network. When it crashes, the leader election module will detect this and elect a new leader.

## Implementation
 - View the implementation [here](./faulty.go).
 - View the implementation of the leader election node [here](../../lib/le.go).
 - View the implementation of the perfect failure detector node [here](../../lib/pfd.go).
 - View the simulation [here](./main.go)