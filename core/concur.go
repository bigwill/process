package core

import (
	"log"
)

func RunChain(src Source, procs []Processor, snk Sink) ([]ControlChannel, MonitorChannel) {
	ctrls := make([]ControlChannel, 2+len(procs), 2+len(procs))
	ctrls[0] = make(ControlChannel)
	ctrls[len(ctrls)-1] = make(ControlChannel)

	monChan := make(MonitorChannel)

	srcOut, snkIn, procCtrls := runProcessors(procs, monChan)

	for i, pc := range procCtrls {
		ctrls[i+1] = pc
	}

	go newSinkRoutine(snk)(snkIn, ctrls[0], monChan)
	go newSourceRoutine(src)(srcOut, ctrls[len(ctrls)-1], monChan)

	return ctrls, monChan
}

func runProcessors(procs []Processor, mon MonitorChannel) (SampleChannel, SampleChannel, []ControlChannel) {
	ctrls := make([]ControlChannel, len(procs), len(procs))
	chainInChan := make(SampleChannel)
	inChan := chainInChan
	var outChan SampleChannel

	for i, proc := range procs {
		ctrls[i] = make(ControlChannel)
		outChan = make(SampleChannel)
		go newProcessorRoutine(proc)(inChan, outChan, ctrls[i], mon)
		inChan = outChan
	}

	return chainInChan, outChan, ctrls
}

func newSourceRoutine(src Source) SourceRoutine {
	return func(out SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
		for {
			select {
			case ctrlVal := <-ctrl:
				if ctrlVal == Quit {
					return
				}
			default:
				v, err := src.Output()
				if err != nil {
					log.Println(err)
					mon <- Quit
				}
				out <- v
			}
		}
	}
}

func newSinkRoutine(snk Sink) SinkRoutine {
	return func(in SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
		for {
			select {
			case ctrlVal := <-ctrl:
				if ctrlVal == Quit {
					return
				}
			case v := <-in:
				err := snk.Input(v)
				if err != nil {
					log.Println(err)
					mon <- Quit
				}
			}
		}
	}
}

func newProcessorRoutine(p Processor) ProcessorRoutine {
	return func(in SampleChannel, out SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
		for {
			select {
			case ctrlVal := <-ctrl:
				if ctrlVal == Quit {
					return
				}
			case v := <-in:
				w, err := p.Process(v)
				if err != nil {
					log.Println(err)
					mon <- Quit
				}
				out <- w
			}
		}
	}
}
