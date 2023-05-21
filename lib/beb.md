# Best-effort Broadcast

## Module
 - Name: BestEffortBroadcast
 - Instance: _beb_

## Messages
 - `p -> beb: BebBroadcast(m)`: Requests to broadcast a message _m_ from a process _p_ to all processes.
 - `beb -> p: BebDeliver(m)`: Indicates that a message _m_ broadcast by process _p_ has been delivered.

## Properties
 - **Validity**: If a correct process broadcasts a message _m_, then every correct process eventually delivers _m_.
 - **No duplication**: No message is delivered more than once.
 - **No creation**: If a process delivers a message _m_, then _m_ was previously broadcast by some process.

## Implementation
 - View the implementation [here](./beb.go).

## Note
This module is not intended to be used when implementing other modules because the disse library has built in a `BroadcastMessage` function that implements broadcasting messages. This module is only included for completeness.