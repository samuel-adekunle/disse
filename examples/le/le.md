# Leader election

## Module
 - Name: LeaderElection
 - Instance: _le_

## Messages
 - `le -> q: LeLeader(p)`: Indicates that a process _p_ is elected as leader.

## Properties
 - **Eventual detection**:  Either there is no correct process, or some correct process is eventually elected as the leader.
 - **Accuracy**:  If a process is leader, then all previously elected leaders have crashed.

## Implementation
 - View the implementation [here](./le.go).
 - View the simulation [here](./main.go).