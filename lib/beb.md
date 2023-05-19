# Best-effort Broadcast

## Module
 - Name: BestEffortBroadcast
 - Instance: _beb_

## Messages
 - `_p_ -> _beb_: BebBroadcast(_m_)`: Requests to broadcast a message _m_ from a process _p_ to all processes.
 - `_beb_ -> _p_: BebDeliver(_m_)`: Indicates that a message _m_ broadcast by process _p_ has been delivered.

## Properties
 - **Validity**: If a correct process broadcasts a message _m_, then every correct process eventually delivers _m_.
 - **No duplication**: No message is delivered more than once.
 - **No creation**: If a process delivers a message _m_, then _m_ was previously broadcast by some process.