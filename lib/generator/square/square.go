package square

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
)

type State struct {
	sampleRate core.Quantity
	f_g core.Param
	T int64 // wave period in terms of samples @ the given sample rate
	i int64 // current index in wave period
}

func MakeState(sampleRate core.Quantity) *State {
	s := &State{sampleRate: sampleRate,
		f_g: linear.MakeState("Freq", "Hz", 30, 10000, .5)}
	s.f_g.SetHandler(func(p core.Param) {
		s.setFrequency(p.Val())
	})
	s.f_g.SetPos(.5)
	return s
}

func (s *State) GetNumParams() core.ParamIdx {
	return 1
}

func (s *State) GetParam(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return s.f_g
	} else {
		return nil
	}
}

func (s *State) Generate() core.Quantity {
	defer func (t *State) {
		t.i = (t.i + 1) % t.T
	}(s)
	if s.i <= s.T / 2 {
		return 1.0
	} else {
		return -1.0
	}
}

func (s *State) setFrequency(f_g core.Quantity) {
	s.T = int64(s.sampleRate / f_g)
}