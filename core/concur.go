package core

func RunChain(src Source, procs []Processor, snk Sink) []ControlChannel {
	ctrls := make([]ControlChannel, 2+len(procs), 2+len(procs))
	ctrls[0] = make(ControlChannel)
	ctrls[len(ctrls)-1] = make(ControlChannel)

	srcOut, snkIn, procCtrls := runProcessors(procs)

	for i, pc := range procCtrls {
		ctrls[i+1] = pc
	}

	go newSinkRoutine(snk)(snkIn, ctrls[len(ctrls)-1])
	go newSourceRoutine(src)(srcOut, ctrls[0])

	return ctrls
}

func runProcessors(procs []Processor) (SampleChannel, SampleChannel, []ControlChannel) {
	ctrls := make([]ControlChannel, len(procs), len(procs))
	chainInChan := make(SampleChannel)
	inChan := chainInChan
	var outChan SampleChannel

	for i, proc := range procs {
		ctrls[i] = make(ControlChannel)
		outChan = make(SampleChannel)
		go newProcessorRoutine(proc)(inChan, outChan, ctrls[i])
		inChan = outChan
	}

	return chainInChan, outChan, ctrls
}

func newSourceRoutine(src Source) SourceRoutine {
	return func(out SampleChannel, ctrl ControlChannel) {
		for {
			select {
			case ctrlVal := <-ctrl:
				if ctrlVal == Quit {
					return
				}
			case out <- src.Output():
			}
		}
	}
}

func newSinkRoutine(snk Sink) SinkRoutine {
	return func(in SampleChannel, ctrl ControlChannel) {
		for {
			select {
			case ctrlVal := <-ctrl:
				if ctrlVal == Quit {
					return
				}
			case v := <-in:
				snk.Input(v)
			}
		}
	}
}

func newProcessorRoutine(p Processor) ProcessorRoutine {
	return func(in SampleChannel, out SampleChannel, ctrl ControlChannel) {
		for {
			select {
			case ctrlVal := <-ctrl:
				if ctrlVal == Quit {
					return
				}
			case v := <-in:
				out <- p.Process(v)
			}
		}
	}
}
