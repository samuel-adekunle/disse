# Perfect point-to-point links

## Module
 - Name: PerfectPointToPointLinks
 - Instance: _pl_

## Messages

 - `a -> pl: PlSend(b, m)`: Requests _pl_ to send a message _m_ from node _a_ to node _b_.
 - `pl -> b: PlDeliver(a, m)`: Delivers message _m_ sent by node _a_ to node _b_.

## Properties

 - **Reliable delivery**: If a correct node _a_ sends a message _m_ to a correct node _b_, then _b_ eventually delivers _m_.
 - **No duplication**: No message is delivered by a node more than once.
 - **No creation**: If some node _b_ delivers a message _m_ with sender _a_, then _m_ was previously sent to _b_ by node _a_.

## Implementation

 - View the implementation [here](./pl.go).