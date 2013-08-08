package main

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
	"github.com/bigwill/process/lib/sink/stdout"
	"github.com/bigwill/process/lib/source/stdin"
)

const bufferSize = 500
const sampleRate = 48000.0

func main() {
	// Source
	stdinS := stdin.NewSource(sampleRate)

	// Processors
	rcfP := rcfilter.NewProcessor(sampleRate)
	rcfP.Param(0).SetPos(.2)

	gainP := gain.NewProcessor(sampleRate)
	gainP.Param(0).SetPos(.1)

	// Sink
	stdoutSink := stdout.NewSink(sampleRate)

	_, monChan := core.RunChain(stdinS, []core.Processor{rcfP, gainP}, stdoutSink)

	for {
		m := <-monChan
		if m.Code == core.Quit {
			return
		}
	}
}