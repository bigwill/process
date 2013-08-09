package main

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/filter"
	"github.com/bigwill/process/lib/sink/play"
	"github.com/bigwill/process/lib/source/square"
	"log"
)

const bufferSize = 500
const sampleRate = 48000.0

func main() {
	// Source
	src := square.NewSource(sampleRate)
	src.Param(0).SetPos(.1)

	// Processors
	filterP := filter.NewProcessor(sampleRate)
	filterP.Param(1).SetPos(.1)
	filterP.Param(0).SetPos(.9)
	filterP.Param(2).SetPos(.9)

	gainP := gain.NewProcessor(sampleRate)
	gainP.Param(0).SetPos(.5)

	// Sink
	snk := play.NewSink(sampleRate)

	// MidiSource
	var midiSrc core.MidiSource = nil

	log.Printf("filter type = %v", filterP.Param(0).Val())
	log.Printf("filter cutoff = %v", filterP.Param(1).Val())
	log.Printf("filter Q = %v", filterP.Param(2).Val())

	_, monChan := core.RunChain(src, []core.Processor{filterP, gainP}, snk, midiSrc)

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