# Best-effort Broadcast

## Module
 - Name: BestEffortBroadcast
 - Instance: _beb_

## Messages
 - `a -> beb: BebBroadcast(m)`: Requests to broadcast a message _m_ from a node _a_ to all nodes.
 - `beb -> b: BebDeliver(a, m)`: Delivers a broadcast message _m_ sent by node _a_ to node _b_.

## Properties
 - **Validity**: If a correct node broadcasts a message _m_, then every correct node eventually delivers _m_.
 - **No duplication**: No message is delivered more than once.
 - **No creation**: If a node delivers a message _m_, then _m_ was previously broadcast by some node.

## Implementation
 - View the implementation [here](./beb.go).