package main

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
	"github.com/bigwill/process/lib/sink/stdout"
	"github.com/bigwill/process/lib/source/stdin"
)

const bufferSize = 500

func main() {
	// Source
	stdinS := stdin.NewSource(48000.0)

	// Processors
	rcfP := rcfilter.NewProcessor(48000.0)
	rcfP.Param(0).SetPos(.2)

	gainP := gain.NewProcessor(48000.0)
	gainP.Param(0).SetPos(.1)

	// Sink
	stdoutSink := stdout.NewSink(48000.0)

	_, monChan := core.RunChain(stdinS, []core.Processor{rcfP, gainP}, stdoutSink)

	for {
		m := <-monChan
		if m.Code == core.Quit {
			return
		}
	}
}