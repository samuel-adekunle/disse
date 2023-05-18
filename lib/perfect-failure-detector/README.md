# Perfect failure detector

## Module
 - Name: PerfectFailureDetector
 - Instance: _pfd_

## Messages
 - `_pfd_ -> p: PfdCrash(_q_)`: Indicates that a process _q_ has crashed.

## Properties
 - **Strong completeness**: Eventually, every process that crashes is permanently detected by every correct process.
 - **Strong accuracy**: If a process p is detected by any process, then p has crashed.