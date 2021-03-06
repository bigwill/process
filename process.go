package main

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/processor/filter"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/sink/play"
	"github.com/bigwill/process/lib/source/square"
	"log"
)

const sampleRate = 48000.0
const numChannels = 2
const framePoolSize = 1000

func main() {
	var err error
	var src core.Source
	var filterP, gainP core.Processor
	var snk core.Sink

	ctx := core.NewContext(sampleRate, numChannels, framePoolSize)

	// Source
	src, err = square.NewSource(ctx)
	if err != nil {
		log.Printf("err = %v", err)
	}

	src.Param(0).SetPos(.1)

	// Processors
	filterP, err = filter.NewProcessor(ctx)
	if err != nil {
		log.Printf("err = %v", err)
	}

	filterP.Param(1).SetPos(.1)
	filterP.Param(0).SetPos(.9)
	filterP.Param(2).SetPos(.9)

	gainP, err = gain.NewProcessor(ctx)
	if err != nil {
		log.Printf("err = %v", err)
	}

	gainP.Param(0).SetPos(.5)

	// Sink
	snk, err = play.NewSink(ctx)
	if err != nil {
		log.Printf("err = %v", err)
	}

	// MidiSource
	var midiSrc core.MidiSource = nil

	log.Printf("filter type = %v", filterP.Param(0).Val())
	log.Printf("filter cutoff = %v", filterP.Param(1).Val())
	log.Printf("filter Q = %v", filterP.Param(2).Val())

	_, monChan := core.RunChain(ctx, src, []core.Processor{filterP, gainP}, snk, midiSrc)

	for {
		m := <-monChan
		switch m := m.(type) {
		case core.ErrorMonitorMessage:
			log.Printf("err = %v", m)
		default:
			log.Printf("mon = %v", m)
		}
	}
}
