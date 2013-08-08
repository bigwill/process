package core

func RunChain(src Source, procs []Processor, snk Sink, midiSrc MidiSource) ([]ControlChannel, MonitorChannel) {
	ctrls := make([]ControlChannel, 2+len(procs), 2+len(procs))
	ctrls[0] = make(ControlChannel)
	ctrls[len(ctrls)-1] = make(ControlChannel)

	monChan := make(MonitorChannel)

	srcOut, snkIn, procCtrls := runProcessors(procs, monChan)

	for i, pc := range procCtrls {
		ctrls[i+1] = pc
	}

	go sinkRoutine(snk, snkIn, ctrls[0], monChan)
	go sourceRoutine(src, srcOut, ctrls[len(ctrls)-1], monChan)

	if midiSrc != nil {
		go midiRoutine(midiSrc, ctrls, monChan)
	}

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
		go processorRoutine(proc, inChan, outChan, ctrls[i], mon)
		inChan = outChan
	}

	return chainInChan, outChan, ctrls
}

func sourceRoutine(src Source, out SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
	for {
		select {
		case ctrlMsg := <-ctrl:
			if ctrlMsg.Code() == Quit {
				return
			}
		default:
			v, err := src.Output()
			if err != nil {
				mon <- MonitorError(src.Name(), err)
			}
			out <- v
		}
	}
}

func sinkRoutine(snk Sink, in SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
	for {
		select {
		case ctrlMsg := <-ctrl:
			if ctrlMsg.Code() == Quit {
				return
			}
		case v := <-in:
			err := snk.Input(v)
			if err != nil {
				mon <- MonitorError(snk.Name(), err)
			}
		}
	}
}

func processorRoutine(p Processor, in SampleChannel, out SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
	for {
		select {
		case ctrlMsg := <-ctrl:
			if ctrlMsg.Code() == Quit {
				return
			}
		case v := <-in:
			w, err := p.Process(v)
			if err != nil {
				mon <- MonitorError(p.Name(), err)
			}
			out <- w
		}
	}
}

func midiRoutine(midiSrc MidiSource, ctrls []ControlChannel, mon MonitorChannel) {
	for {
		midiMsg, err := midiSrc.Output()
		if err != nil {
			mon <- MonitorError(midiSrc.Name(), err)
		}

		if int(midiMsg.Channel()) < len(ctrls) {
			ctrls[midiMsg.Channel()] <- midiSrc.Mapper().Map(midiMsg)
		}
	}
}
