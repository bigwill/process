package square

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
	"github.com/bigwill/process/lib/processor/filter"
)

type State struct {
	sampleRate core.Quantity
	f_g        core.Param
	t          int64          // wave period in terms of samples @ the given sample rate
	i          int64          // current index in wave period
	f_a1       core.Processor // 2 anti-aliasing filters for 24dB rolloff around 20 kHZ
	f_a2       core.Processor
}

func NewSource(sampleRate core.Quantity) core.Source {
	s := &State{sampleRate: sampleRate,
		f_g:  linear.NewState("Freq", "Hz", 30, 10000, .5),
		f_a1: filter.NewProcessor(sampleRate),
		f_a2: filter.NewProcessor(sampleRate)}
	s.f_g.SetHandler(func(p core.Param) {
		s.setFrequency(p.Val())
	})
	s.f_g.SetPos(.5)

	s.f_a1.Param(1).SetPos(1)   // cutoff ~= 20kHz
	s.f_a1.Param(2).SetPos(.04) // Q ~= .72
	s.f_a2.Param(1).SetPos(1)
	s.f_a2.Param(2).SetPos(.04)

	return s
}

func (s *State) Name() string {
	return "Sq Osc"
}

func (s *State) NumParams() core.ParamIdx {
	return 1
}

func (s *State) Param(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return s.f_g
	} else {
		return nil
	}
}

func (s *State) squareOutput() core.Quantity {
	defer func(q *State) {
		q.i = (q.i + 1) % q.t
	}(s)
	if s.i <= s.t/2 {
		return 1.0
	} else {
		return -1.0
	}
}

func (s *State) Output() (core.Quantity, error) {
	y, err := s.f_a1.Process(s.squareOutput())
	if err != nil {
		return 0, err
	}
	return s.f_a2.Process(y)
}

func (s *State) setFrequency(f_g core.Quantity) {
	s.t = int64(s.sampleRate / f_g)
}
