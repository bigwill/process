package gain

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
)

type state struct {
	g core.Param
}

func NewProcessor(ctx core.Context) (core.Processor, error) {
	return &state{g: linear.NewState("Gain", "dB", 0, 1, .8)}, nil
}

func (s *state) Name() string {
	return "Gain"
}

func (s *state) NumParams() core.ParamIdx {
	return 1
}

func (s *state) Param(idx core.ParamIdx) core.Param {
	if idx == 0 {
		return s.g
	} else {
		return nil
	}
}

func (s *state) Process(x core.SampleFrame) error {
	for i := core.Integer(0); i < x.NumChannels(); i++ {
		x.SetChannelVal(i, s.g.Val()*x.ChannelVal(i))
	}

	return nil
}
