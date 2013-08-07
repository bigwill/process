package main

import (
	"encoding/binary"
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/generator/square"
	"github.com/bigwill/process/lib/processor/gain"
	"github.com/bigwill/process/lib/processor/rcfilter"
	"log"
	"os"
)

const bufferSize = 500

func main() {
	sqG := square.NewGenerator(48000.0)
	sqG.Param(0).SetPos(.1)

	rcfP := rcfilter.NewProcessor(48000.0)
	rcfP.Param(0).SetPos(.2)

	gainP := gain.NewProcessor(48000.0)
	gainP.Param(0).SetPos(.1)
	log.Printf("%v", gainP.Param(0).Val())

	outChan, _ := core.RunGeneratorChain(sqG, []core.Processor{rcfP, gainP})

	log.Printf("%s %s", sqG.Name(), sqG.Param(0).Repr())
	log.Printf("%s %s", rcfP.Name(), rcfP.Param(0).Repr())

	buf := make([]core.Quantity, bufferSize, bufferSize)
	for {
		for i := 0; i < bufferSize; i++ {
			s := <-outChan
			buf[i] = s
		}
		err := binary.Write(os.Stdout, binary.LittleEndian, buf)
		if err != nil {
			log.Printf("%v", err)
			return
		}
	}
}
