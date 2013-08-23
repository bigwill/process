package gain

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
)

type State struct {
	g core.Param
}

func NewProcessor(ctx core.Context) (core.Processor, error) {
	return &State{g: linear.NewState("Gain", "dB", 0, 1, .8)}, nil
}

func (s *State) Name() string {
	return "Gain"
}

func (s *State) NumParams() core.ParamIdx {
	return 1
}

func (s *State) Param(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return s.g
	} else {
		return nil
	}
}

func (s *State) Process(x core.SampleFrame) (core.SampleFrame, error) {
	for i := core.Index(0); i < x.NumChannels(); i++ {
		x.SetChannelVal(i, s.g.Val()*x.ChannelVal(i))
	}

	return x, nil
}
