package main

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
	"github.com/bigwill/process/lib/sink/stdout"
	"github.com/bigwill/process/lib/source/square"
	"log"
)

const bufferSize = 500

func main() {
	// Source
	sqG := square.NewSource(48000.0)
	sqG.Param(0).SetPos(.1)

	// Processors
	rcfP := rcfilter.NewProcessor(48000.0)
	rcfP.Param(0).SetPos(.2)

	gainP := gain.NewProcessor(48000.0)
	gainP.Param(0).SetPos(.1)

	// Sink
	stdoutSink := stdout.NewSink(48000.0)

	_ = core.RunChain(sqG, []core.Processor{rcfP, gainP}, stdoutSink)

	log.Printf("%s %s", sqG.Name(), sqG.Param(0).Repr())
	log.Printf("%s %s", rcfP.Name(), rcfP.Param(0).Repr())

	quit := make(chan int)
	<-quit
}