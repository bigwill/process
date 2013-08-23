package rcfilter

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
	"math"
)

type processorState struct {
	ctx core.ProcessorContext
	f_c core.Param
	k   core.Quantity
	l   core.Quantity

	channels []*channelState
}

type channelState struct {
	v_o_n1 core.Quantity
	v_i_n1 core.Quantity
	i_n1   core.Quantity
}

func NewProcessor(ctx core.ProcessorContext) core.Processor {
	// TODO: l value is zero by default for now
	s := &processorState{ctx: ctx,
		f_c: linear.NewState("Cutoff", "Hz", 100.0, 10000.0, .5)}

	s.channels = make([]*channelState, ctx.NumChannels())
	for i := core.Index(0); i < ctx.NumChannels(); i++ {
		s.channels[i] = &channelState{}
	}

	s.f_c.SetHandler(func(p core.Param) {
		s.setCutoff(p.Val())
	})
	s.f_c.SetPos(.5)
	return s
}

func (ps *processorState) Name() string {
	return "RC Fil"
}

func (ps *processorState) NumParams() core.ParamIdx {
	return 1
}

func (ps *processorState) Param(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return ps.f_c
	} else {
		return nil
	}
}

func (ps *processorState) Process(fr core.SampleFrame) (core.SampleFrame, error) {
	for i := core.Index(0); i < fr.NumChannels(); i++ {
		fr.SetChannelVal(i, ps.channels[i].process(ps, fr.ChannelVal(i)))
	}
	return fr, nil
}

func (cs *channelState) process(ps *processorState, v_i_n core.Quantity) core.Quantity {
	// Processing
	i_n := cs.v_i_n1 - cs.v_o_n1
	v_o_n := cs.v_o_n1 + ps.k*i_n + ps.l*(i_n-cs.i_n1)

	// Update processorState for next sample
	cs.v_i_n1 = v_i_n
	cs.v_o_n1 = v_o_n
	cs.i_n1 = i_n

	return v_o_n
}

func (ps *processorState) setCutoff(f_c core.Quantity) {
	var RC core.Quantity = 1 / (2.0 * math.Pi * f_c)
	ps.k = 1 / (ps.ctx.SampleRate() * RC)
}
