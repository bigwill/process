package stdin

import (
	"encoding/binary"
	"github.com/bigwill/process/core"
	"os"
)

const bufferSize = 500

type state struct {
	ctx core.Context
	i   int
	buf []core.Quantity
}

func NewSource(ctx core.Context) core.Source {
	return &state{ctx: ctx, buf: make([]core.Quantity, bufferSize*ctx.NumChannels())}
}

func (s *state) Name() string {
	return "Std In"
}

func (s *state) NumParams() core.ParamIdx {
	return 0
}

func (s *state) Param(idx core.ParamIdx) core.Param {
	return nil
}

func (s *state) Output() (core.SampleFrame, error) {
	if s.i == len(s.buf) {
		err := binary.Read(os.Stdin, binary.LittleEndian, s.buf)
		if err != nil {
			return nil, err
		}

		s.i = 0
	}

	fr := s.ctx.FramePool().DequeueFrame()
	for j := core.Index(0); j < s.ctx.NumChannels(); j++ {
		fr.SetChannelVal(j, s.buf[s.i])
		s.i++
	}

	return fr, nil
}
