# Using process

After building the project, just run the process command.

Requires sox and play commands available in the sox project (http://sox.sourceforge.net/)

Needs a lot of work to be more interesting

# To-dos

- introduce utils for common source, sink, processor patterns
- maybe base implementations for sources, sinks, processors?
- do some profiling
- frame pool size computation
- parallel chains
- more processors!
- revisit cgo at some point
- support multiple sources and sinks
- midisource -> controlsource (introduce a cv-like control model?)
- additional parameter types: logarithmic, exponential