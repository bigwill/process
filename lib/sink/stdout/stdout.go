package stdout

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

func NewSink(ctx core.Context) core.Sink {
	return &state{ctx: ctx, buf: make([]core.Quantity, bufferSize*ctx.NumChannels())}
}

func (s *state) Name() string {
	return "Std Out"
}

func (s *state) NumParams() core.ParamIdx {
	return 0
}

func (s *state) Param(idx core.ParamIdx) core.Param {
	return nil
}

func (s *state) Input(fr core.SampleFrame) error {
	for j := core.Index(0); j < fr.NumChannels(); j++ {
		s.buf[s.i] = fr.ChannelVal(j)
		s.i++
	}

	s.ctx.FramePool().EnqueueFrame(fr)

	if s.i == len(s.buf) {
		err := binary.Write(os.Stdout, binary.LittleEndian, s.buf)
		if err != nil {
			return err
		}

		s.i = 0
	}

	return nil
}
