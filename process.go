package main

import (
	"fmt"
	"log"
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/generator/square"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
)

func main() {
	var square core.Generator = square.MakeGenerator(48000.0)
	var filter core.Processor = rcfilter.MakeProcessor(48000.0)
	var gain core.Processor = gain.MakeProcessor(48000.0)
	filter.Param(0).SetPos(.2)
	square.Param(0).SetPos(.6)

	log.Printf("%s %s", square.Name(), square.Param(0).Repr())
	log.Printf("%s %s", filter.Name(), filter.Param(0).Repr())

	for i := 0; i < 4800; i++ {
		fmt.Printf("%v\t%v\n", i, gain.Process(filter.Process(square.Generate())))
	}
}
