package main

import (
	"encoding/binary"
	"log"
	"os"
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/generator/square"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
)

const bufferSize = 500

func main() {
	var square core.Generator = square.NewGenerator(48000.0)
	var filter core.Processor = rcfilter.NewProcessor(48000.0)
	var gain core.Processor = gain.NewProcessor(48000.0)
	filter.Param(0).SetPos(.2)
	square.Param(0).SetPos(.1)
	gain.Param(0).SetPos(.5)

	log.Printf("%s %s", square.Name(), square.Param(0).Repr())
	log.Printf("%s %s", filter.Name(), filter.Param(0).Repr())

	buf := make([]core.Quantity, bufferSize, bufferSize)
	for {
		for i := 0; i < bufferSize; i++ {
			buf[i] = gain.Process(filter.Process(square.Generate()))
		}
		err := binary.Write(os.Stdout, binary.LittleEndian, buf)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}
