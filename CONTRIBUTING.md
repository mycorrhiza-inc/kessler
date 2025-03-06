# Licensing

Currently kessler, as well as all previous versions are under an MIT License, however we might change the license to another with more copyleft protections.

Whatever we do with it, we will do the following.

- The all source code for kessler (both backend and frotnend )will be publicly available for all future versions. 

- Self hosting kessler, for non-profits, small teams and personal use will always be free of charge. And all of the backend code needed to run an open source version will always be availible.

- Contributors will never need to sign a CLA and will always retain copyright of their code.

# Code Advice Stuff
## Naming Conventions

- Order matters for readability (even if it doesn't affect semantics). On the first read, a file is read top-down, so put important things near the top. The main function goes first.

-  Add units or qualifiers to variable names, and put the units or qualifiers last, sorted by descending significance, so that the variable starts with the most significant word, and ends with the least significant word. For example, latency_ms_max rather than max_latency_ms. This will then line up nicely when latency_ms_min is added, as well as group all variables that relate to latency.

## Principles stolen from tigerbeatle:


- Put a limit on everything because, in reality, this is what we expect—everything has a limit. For example, all loops and all queues must have a fixed upper bound to prevent infinite loops or tail latency spikes. This follows the “fail-fast” principle so that violations are detected sooner rather than later. Where a loop cannot terminate (e.g. an event loop), this must be asserted.


- The golden rule of assertions is to assert the positive space that you do expect AND to assert the negative space that you do not expect because where data moves across the valid/invalid boundary between these spaces is where interesting bugs are often found. This is also why tests must test exhaustively, not only with valid data but also with invalid data, and as valid data becomes invalid.

- Assertions detect programmer errors. Unlike operating errors, which are expected and which must be handled, assertion failures are unexpected. The only correct way to handle corrupt code is to crash. Assertions downgrade catastrophic correctness bugs into liveness bugs. 

- On occasion, you may use a blatantly true assertion instead of a comment as stronger documentation where the assertion condition is critical and surprising.


### Potentially good stuff:

- Use only very simple, explicit control flow for clarity. Do not use recursion to ensure that all executions that should be bounded are bounded. Use only a minimum of excellent abstractions but only if they make the best sense of the domain. Abstractions are never zero cost. Every abstraction introduces the risk of a leaky abstraction.


