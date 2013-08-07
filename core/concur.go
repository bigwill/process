package core

func RunGeneratorChain(gen Generator, procs []Processor) (SampleChannel, []ControlChannel) {
	ctrls := make([]ControlChannel, 1+len(procs), 1+len(procs))
	genOut, outChan, procCtrls := RunProcessorChain(procs)

	for i, pc := range procCtrls {
		ctrls[i+1] = pc
	}

	go newGeneratorRoutine(gen)(genOut, ctrls[0])
	return outChan, ctrls
}

func RunProcessorChain(procs []Processor) (SampleChannel, SampleChannel, []ControlChannel) {
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

func newGeneratorRoutine(g Generator) GeneratorRoutine {
	return func(out SampleChannel, ctrl ControlChannel) {
		for {
			select {
			case ctrlVal := <-ctrl:
				if ctrlVal == Quit {
					return
				}
			case out <- g.Generate():
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
			case out <- p.Process(<-in):
			}
		}
	}
}
