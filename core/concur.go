package core

func RunChain(gen Generator, procs []Processor) (chan Quantity, []chan Control) {
	ctrls := make([]chan Control, 1+len(procs), 1+len(procs))
	var inChan, outChan chan Quantity

	ctrls[0] = make(chan Control)
	genOut := make(chan Quantity)
	go newGeneratorRoutine(Generator(gen), genOut, ctrls[0])()
	inChan = genOut

	for i, proc := range procs {
		ctrls[i+1] = make(chan Control)
		outChan = make(chan Quantity)
		go newProcessorRoutine(Processor(proc), inChan, outChan, ctrls[i+1])()
		inChan = outChan
	}

	return outChan, ctrls
}

func newGeneratorRoutine(g Generator, out chan Quantity, ctrl chan Control) func() {
	return func() {
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

func newProcessorRoutine(p Processor, in, out chan Quantity, ctrl chan Control) func() {
	return func() {
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
