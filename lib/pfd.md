# Perfect failure detector

## Module
 - Name: PerfectFailureDetector
 - Instance: _pfd_

## Messages
 - `pfd -> a: PfdCrash(b)`: Indicates that a node _b_ has crashed.
 - `pfd -> a: HeartbeatRequest`: Requests a heartbeat reply from a node _a_.
 - `a -> pfd: HeartbeatReply`: Replies to a heartbeat request from _pfd_.

## Properties
 - **Strong completeness**: Eventually, every node that crashes is permanently detected by every correct node.
 - **Strong accuracy**: If a node _a_ is detected by any node, then _a_ has crashed.

## Implementation
 - View the implementation [here](./pfd.go).