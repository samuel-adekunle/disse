# List of TODOs

## Implemented features

- [X] Node creation
- [X] Message passing
- [X] Timers and timing structures
- [X] Message Latency
- [X] Logging
- [X] Ping Pong example
- [X] Node composition (ability to reuse code between nodes)
- [X] Sequence Diagram Visualisation
  - [X] Using [PlantUML](https://plantuml.com/en-dark/sequence-diagram)
- [X] Unit testing framework for the system
- [X] Interrupts: Sleep, Start, Stop (Permanent)
- [X] Standardize dropping messages
- [X] Finish adding examples as documentation of library usage (especially tests)

## Pending work

- [ ] Make logging DRY, use a centralised logging interface or function.
- [ ] Examples
  - [ ] Perfect Point to Point Links
  - [ ] Reliable Broadcast using Best Effort Broadcast which uses Perfect Point to Point Links.
- [ ] Create detailed project plan
  - [ ] Find students to test out the tool during summer exams
  - [ ] Feedback collection (More Quantitative vs Qualitative feedback)
  - [ ] Create evaluation plan
    - [ ] Ability to discover faults in existing systems
    - [ ] Speed and writing and deploying new systems
    - [ ] Observe users using the tool
  - [ ] Start writing report

## Extensions

- [ ] Busy nodes? Consider only handling messages when the node is not busy i.e. Need to keep track of busy nodes
- [ ] Introduce more network effects like packet loss, packet corruption, packet reordering, jitter
- [ ] Add a Restart Stopped Nodes Interrupt
- [ ] Implement a custom timer to allow pausing, speeding up and slowing down the simulation
- [ ] Integration / Runtime testing: i.e. invariants like no duplication, no timeout while waiting for response
- [ ] Interactivity: requires a custom timer that can be paused, sped up or slowed down
  - [ ] Also allow faults to be injected at runtime (can be implemented by having a node that listens to stdin and sends messages in system) i.e. UserNode
- [ ] Model checking with state pruning
- [ ] Simulator Interface, with extension interfaces like a TCP / UDP simulator (needs more complex types)
  - [ ] Enable to run actual distributed systems using cloud services and test different protocols
- [ ] Specify system using a configuration language like YAML or JSON
- [ ] More examples
  - [ ] Paxos
  - [ ] Bitcoin Blockchain
  - [ ] Shared store example
