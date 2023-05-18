# Perfect point-to-point links

## Module
 - Name: PerfectPointToPointLinks
 - Instance: _pl_

## Messages

 - `_p_ -> _pl_: PlSend(_q_, _m_)`: Requests to send a message _m_ from process _p_ to process _q_.
 - `_pl_ -> _p_: PlDeliver(_m_)`: Indicates that message _m_ sent by process _p_ has been delivered.  

## Properties

 - **Reliable delivery**: If a correct process _p_ sends a message _m_ to a correct process _q_, then _q_ eventually delivers _m_.
 - **No duplication**: No message is delivered by a process more than once.
 - **No creation**: If some process _q_ delivers a message _m_ with sender _p_, then _m_ was previously sent to _q_ by process _p_.

## Implementation

 - View the implementation [here](./pl.go).