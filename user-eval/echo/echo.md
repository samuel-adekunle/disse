# Echo

## Module
 - Name: Echo
 - Instance: _echo_

## Messages
 - `p -> echo: EchoSend(m)`: Requests that a message _m_ sent by process _p_ is sent back to process _p_.
 - `echo -> p: EchoDeliver(m)`: Indicates that a message _m_ sent by process _p_ has been sent back.

## Properties
 - **Echo**: If a correct process _p_ sends a message _m_ to the echo server, then _p_ eventually delivers _m_ back from the echo server.