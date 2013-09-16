package core

const controlChannelBufferSize = 100

func RunChain(ctx Context, src Source, procs []Processor, snk Sink, midiSrc MidiSource) ([]ControlChannel, MonitorChannel) {
	ctrls := make([]ControlChannel, 2+len(procs), 2+len(procs))
	ctrls[0] = makeBufferControlChannel()
	ctrls[len(ctrls)-1] = makeBufferControlChannel()

	monChan := make(MonitorChannel)

	srcOut, snkIn, procCtrls := runProcessors(ctx, procs, monChan)

	for i, pc := range procCtrls {
		ctrls[i+1] = pc
	}

	go sinkRoutine(ctx, snk, snkIn, ctrls[0], monChan)
	go sourceRoutine(ctx, src, srcOut, ctrls[len(ctrls)-1], monChan)

	if midiSrc != nil {
		go midiRoutine(ctx, midiSrc, ctrls, monChan)
	}

	return ctrls, monChan
}

func makeBufferControlChannel() ControlChannel {
	return make(ControlChannel, controlChannelBufferSize)
}

func runProcessors(ctx Context, procs []Processor, mon MonitorChannel) (SampleChannel, SampleChannel, []ControlChannel) {
	ctrls := make([]ControlChannel, len(procs), len(procs))
	chainInChan := make(SampleChannel)
	inChan := chainInChan
	var outChan SampleChannel

	for i, proc := range procs {
		ctrls[i] = makeBufferControlChannel()
		outChan = make(SampleChannel)
		go processorRoutine(ctx, proc, inChan, outChan, ctrls[i], mon)
		inChan = outChan
	}

	return chainInChan, outChan, ctrls
}

func sourceRoutine(ctx Context, src Source, out SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
	for {
		select {
		case ctrlMsg := <-ctrl:
			if ctrlMsg.Code() == Quit {
				return
			}
		default:
			v := ctx.FramePool().DequeueFrame()
			err := src.Output(v)
			if err != nil {
				mon <- MonitorError(src.Name(), err)
			}
			out <- v
		}
	}
}

func sinkRoutine(ctx Context, snk Sink, in SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
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
			ctx.FramePool().EnqueueFrame(v)
		}
	}
}

func processorRoutine(ctx Context, p Processor, in SampleChannel, out SampleChannel, ctrl ControlChannel, mon MonitorChannel) {
	for {
		select {
		case ctrlMsg := <-ctrl:
			if ctrlMsg.Code() == Quit {
				return
			}
		case v := <-in:
			err := p.Process(v)
			if err != nil {
				mon <- MonitorError(p.Name(), err)
			}
			out <- v
		}
	}
}

func midiRoutine(ctx Context, midiSrc MidiSource, ctrls []ControlChannel, mon MonitorChannel) {
	for {
		midiMsg, err := midiSrc.Output()
		if err != nil {
			mon <- MonitorError(midiSrc.Name(), err)
			continue
		}

		if int(midiMsg.Channel()) < len(ctrls) {
			ctrls[midiMsg.Channel()] <- midiSrc.Mapper().Map(midiMsg)
		}
	}
}
