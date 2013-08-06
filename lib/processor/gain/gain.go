package gain

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
)

type State struct {
	g core.Param
}

func MakeState(sample_rate core.Quantity) *State {
	return &State{g: linear.MakeState("Gain", "dB", 0, 1, .8)}
}

func (s *State) GetNumParams() core.ParamIdx {
	return 1
}

func (s *State) GetParam(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return s.g
	} else {
		return nil
	}
}

func (s *State) Process(v_i_n core.Quantity) core.Quantity {
	return s.g.Val() * v_i_n
}