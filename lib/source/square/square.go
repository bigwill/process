package square

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
	"github.com/bigwill/process/lib/processor/filter"
)

type state struct {
	ctx  core.Context
	f_g  core.Param
	t    int64          // wave period in terms of samples @ the given sample rate
	i    int64          // current index in wave period
	f_a1 core.Processor // 2 anti-aliasing filters for 24dB rolloff around 20 kHZ
	f_a2 core.Processor
}

func NewSource(ctx core.Context) (core.Source, error) {
	var err error
	var f_a1, f_a2 core.Processor

	f_a1, err = filter.NewProcessor(ctx)
	if err != nil {
		return nil, err
	}

	f_a2, err = filter.NewProcessor(ctx)
	if err != nil {
		return nil, err
	}

	s := &state{ctx: ctx,
		f_g:  linear.NewState("Freq", "Hz", 30, 10000, .5),
		f_a1: f_a1,
		f_a2: f_a2}
	s.f_g.SetHandler(func(p core.Param) {
		s.setFrequency(p.Val())
	})
	s.f_g.SetPos(.5)

	s.f_a1.Param(1).SetPos(1)   // cutoff ~= 20kHz
	s.f_a1.Param(2).SetPos(.04) // Q ~= .72
	s.f_a2.Param(1).SetPos(1)
	s.f_a2.Param(2).SetPos(.04)

	return s, nil
}

func (s *state) Name() string {
	return "Sq Osc"
}

func (s *state) NumParams() core.ParamIdx {
	return 1
}

func (s *state) Param(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return s.f_g
	} else {
		return nil
	}
}

func (s *state) squareOutput() core.Quantity {
	defer func(q *state) {
		q.i = (q.i + 1) % q.t
	}(s)
	if s.i <= s.t/2 {
		return 1.0
	} else {
		return -1.0
	}
}

func (s *state) Output() (core.SampleFrame, error) {
	fr := s.ctx.FramePool().DequeueFrame()

	v := s.squareOutput()
	for i := core.Index(0); i < s.ctx.NumChannels(); i++ {
		fr.SetChannelVal(i, v)
	}

	fr, err := s.f_a1.Process(fr)
	if err != nil {
		return nil, err
	}

	fr, err = s.f_a2.Process(fr)
	if err != nil {
		return nil, err
	}

	return fr, nil
}

func (s *state) setFrequency(f_g core.Quantity) {
	s.t = int64(s.ctx.SampleRate() / f_g)
}
