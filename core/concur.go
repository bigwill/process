package core

func RunChain(gen Generator, procs []Processor) (SampleChannel, []ControlChannel) {
	ctrls := make([]ControlChannel, 1+len(procs), 1+len(procs))
	var inChan, outChan SampleChannel

	ctrls[0] = make(ControlChannel)
	genOut := make(SampleChannel)
	go newGeneratorRoutine(gen)(genOut, ctrls[0])
	inChan = genOut

	for i, proc := range procs {
		ctrls[i+1] = make(ControlChannel)
		outChan = make(SampleChannel)
		go newProcessorRoutine(proc)(inChan, outChan, ctrls[i+1])
		inChan = outChan
	}

	return outChan, ctrls
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
