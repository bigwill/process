package main

import (
	"fmt"
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/generator/square"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
)

func main() {
	var square core.Generator = square.MakeGenerator(48000.0)
	var filter core.Processor = rcfilter.MakeProcessor(48000.0)
	var gain core.Processor = gain.MakeProcessor(48000.0)
	gain.GetParam(0).SetPos(.2)
	square.GetParam(0).SetPos(.6)

	for i := 0; i < 4800; i++ {
		fmt.Printf("%v\t%v\n", i, gain.Process(filter.Process(square.Generate())))
	}
}
