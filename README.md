# Using process

After building the project, just run the process command.

Requires sox and play commands available in the sox project (http://sox.sourceforge.net/)

Needs a lot of work to be more interesting

# To-dos

- modify Source, Sink, and Processor interfaces and refactor implementers so they don't need to know how to deal with frame pools (dealing with the frame pool will be factored out into the top level goroutine functions in concur.go)
- introduce utils for common source, sink, processor patterns
- maybe base implementations for sources, sinks, processors?
- do some profiling
- frame pool size computation
- parallel chains
- more processors!
- revisit cgo at some point
- support multiple sources and sinks
- introduce opportunity to report errors during source/sink/processor initialization
- midisource -> controlsource (introduce a cv-like control model?)
- additional parameter types: logarithmic, exponential