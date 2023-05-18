# Perfect failure detector

## Module
 - Name: PerfectFailureDetector
 - Instance: _pfd_

## Messages
 - `_pfd_ -> _p_: PfdCrash(_q_)`: Indicates that a process _q_ has crashed.

## Properties
 - **Strong completeness**: Eventually, every process that crashes is permanently detected by every correct process.
 - **Strong accuracy**: If a process _p_ is detected by any process, then _p_ has crashed.

## Implementation
 - View the implementation [here](./pfd.go).