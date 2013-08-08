package main

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/midimapper/basic"
	"github.com/bigwill/process/lib/midisource/ps3midi"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
	"github.com/bigwill/process/lib/sink/play"
	"github.com/bigwill/process/lib/source/sox"
	"log"
)

const bufferSize = 500
const sampleRate = 48000.0

func main() {
	// Source
	src := sox.NewSource(sampleRate)

	// Processors
	rcfP := rcfilter.NewProcessor(sampleRate)
	rcfP.Param(0).SetPos(.2)

	gainP := gain.NewProcessor(sampleRate)
	gainP.Param(0).SetPos(.5)

	// Sink
	snk := play.NewSink(sampleRate)

	// MidiSource
	midiSrc := ps3midi.NewMidiSource(basic.NewMidiMapper())

	_, monChan := core.RunChain(src, []core.Processor{rcfP, gainP}, snk, midiSrc)

	for {
		m := <-monChan
		switch m := m.(type) {
		case core.ErrorMonitorMessage:
			log.Printf("err = %v", m.Err())
			return
		default:
			log.Printf("code = %v", m.Code())
		}
	}
}