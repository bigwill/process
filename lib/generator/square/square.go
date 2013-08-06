package square

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
)

type State struct {
	sampleRate core.Quantity
	f_g core.Param
	t int64 // wave period in terms of samples @ the given sample rate
	i int64 // current index in wave period
}

func NewGenerator(sampleRate core.Quantity) core.Generator {
	s := &State{sampleRate: sampleRate,
		f_g: linear.NewState("Freq", "Hz", 30, 10000, .5)}
	s.f_g.SetHandler(func(p core.Param) {
		s.setFrequency(p.Val())
	})
	s.f_g.SetPos(.5)
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

func (s *State) Generate() core.Quantity {
	defer func (q *State) {
		q.i = (q.i + 1) % q.t
	}(s)
	if s.i <= s.t / 2 {
		return 1.0
	} else {
		return -1.0
	}
}

func (s *State) setFrequency(f_g core.Quantity) {
	s.t = int64(s.sampleRate / f_g)
}