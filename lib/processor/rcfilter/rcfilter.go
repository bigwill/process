package rcfilter

import (
	"math"
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
)

type State struct {
	sampleRate core.Quantity
	f_c core.Param
	v_o_n1 core.Quantity
	v_i_n1 core.Quantity
	i_n1 core.Quantity
	k core.Quantity
	l core.Quantity
}

func MakeProcessor(sampleRate core.Quantity) core.Processor {
	// TODO: l value is zero by default for now
	s := &State{sampleRate: sampleRate,
		f_c: linear.MakeState("Cutoff", "Hz", 100.0, 10000.0, .5)}
	s.f_c.SetHandler(func (p core.Param) {
		s.setCutoff(p.Val())
	})
	s.f_c.SetPos(.5)
	return s
}

func (s *State) NumParams() core.ParamIdx {
	return 1
}

func (s *State) Param(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return s.f_c
	} else {
		return nil
	}
}

func (s *State) Process(v_i_n core.Quantity) core.Quantity {
	// Processing
	i_n := s.v_i_n1 - s.v_o_n1
	v_o_n := s.v_o_n1 + s.k*i_n + s.l * (i_n - s.i_n1)

	// Update state for next sample
	s.v_i_n1 = v_i_n
	s.v_o_n1 = v_o_n
	s.i_n1 = i_n

	return v_o_n
}

func (s *State) setCutoff(f_c core.Quantity) {
	var RC core.Quantity = 1 / (2.0 * math.Pi * f_c)
	s.k = 1 / (s.sampleRate * RC)
}
