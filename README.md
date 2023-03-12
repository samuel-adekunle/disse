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

## Short term features

- [ ] TCP example
  - [ ] File transfer example built on top of TCP
- [ ] Create detailed project plan

## Long term features

- [ ] Paxos example
- [ ] Blockchain example
- [ ] Shared store example
- [ ] Discover faults in existing systems
- Simulator Interface, with extension interfaces like a TCP / UDP simulator
- Fault injection using custom timer that can be paused
- Invariant testing like no duplication / response timeout checked on message / timer queue with current setup
- Will need more complex message, timer and address types (but that was always expected. Implement as an interface t
hat implements encodable or serializable or something, same with timer
- address depends on simulator type, can be simple string for shared memory simulator or ip addresses for tcp/ip 
- run simulator on aws / gcp across different regions rather than injecting fake latency

## Extensions

- [ ] Model checking with state pruning
- [ ] Fault Injection
- [ ] Specify system using a configuration language like YAML or JSON
