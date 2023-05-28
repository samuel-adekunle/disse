# Leader election

## Module
 - Name: LeaderElection
 - Instance: _le_

## Messages
 - `le -> b: LeLeader(a)`: Indicates that a node _a_ is elected as leader.

## Properties
 - **Eventual detection**:  Either there is no correct node, or some correct node is eventually elected as the leader.
 - **Accuracy**:  If a node is leader, then all previously elected leaders have crashed.