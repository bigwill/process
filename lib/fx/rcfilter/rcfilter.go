package rcfilter

import (
	"math"
	"github.com/bigwill/process/core"
)

type State struct {
	v_o_n1 core.Quantity
	v_i_n1 core.Quantity
	i_n1 core.Quantity
	k core.Quantity
	l core.Quantity
}

func MakeState(sample_rate core.Quantity) *State {
	var f_c core.Quantity = 1000.0 // 1kHZ cutoff frequency
	var RC core.Quantity = 1 / (2.0 * math.Pi * f_c)

	// TODO: l value is zero by default for now
	T := 1 / sample_rate
	s := &State{k: T / RC}
	return s
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