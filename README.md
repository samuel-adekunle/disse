# Distributed Systems Discrete Event Simulator (DISSE)

View report [here](https://www.overleaf.com/project/640244e50c059ab2d68e51de).

A lot of inspiration from [DS Labs](https://github.com/emichael/dslabs).

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

## Pending work

- [ ] Fault Injection
  - [ ] Allow certain operations on nodes after certain amounts of time: Sleep, Stop, Start, Restart using a separate queues
- [ ] Busy nodes? Consider only handling messages when the node is not busy i.e. Need to keep track of busy nodes
- [ ] TCP example
  - [ ] File transfer example built on top of TCP
- [ ] Create detailed project plan
  - [ ] Find students to test out the tool during summer exams
  - [ ] Feedback collection (More Quantitative vs Qualitative feedback)
  - [ ] Create evaluation plan
    - [ ] Ability to discover faults in existing systems
    - [ ] Speed and writing and deploying new systems

## Extensions

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
