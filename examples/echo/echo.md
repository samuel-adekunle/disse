# Echo

## Module
 - Name: Echo
 - Instance: _echo_

## Messages
 - `a -> echo: EchoSend(m)`: Requests that a message _m_ sent by node _a_ is sent back to node _a_.
 - `echo -> a: EchoDeliver(m)`: Delivers a message _m_ sent by node _a_ back to node _a_.

## Properties
 - **Echo**: If a correct node _a_ sends a message _m_ to the echo server, then _a_ eventually delivers _m_ back from the echo server.

## Implementation
 - View the implementation [here](./echo.go).
 - View the simulation [here](./main.go).