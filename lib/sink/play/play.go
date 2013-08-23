package play

import (
	"encoding/binary"
	"fmt"
	"github.com/bigwill/process/core"
	"io"
	"os/exec"
)

const bufferSize = 500

type state struct {
	ctx    core.Context
	i      int
	buf    []core.Quantity
	writer io.Writer
}

func NewSink(ctx core.Context) (core.Sink, error) {
	s := &state{ctx: ctx, buf: make([]core.Quantity, bufferSize*ctx.NumChannels())}

	// TODO: yuck. make this cleaner
	cmd := exec.Command("play", "-t", "f64", "-r", fmt.Sprintf("%v", ctx.SampleRate()), "-c", fmt.Sprintf("%v", ctx.NumChannels()), "-")

	var err error
	s.writer, err = cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *state) Name() string {
	return "Sox Play"
}

func (s *state) NumParams() core.ParamIdx {
	return 0
}

func (s *state) Param(idx core.ParamIdx) core.Param {
	return nil
}

func (s *state) Input(v core.SampleFrame) error {
	for j := core.Index(0); j < v.NumChannels(); j++ {
		s.buf[s.i] = v.ChannelVal(j)
		s.i++
	}

	s.ctx.FramePool().EnqueueFrame(v)

	if s.i == len(s.buf) {
		err := binary.Write(s.writer, binary.LittleEndian, s.buf)
		if err != nil {
			return err
		}

		s.i = 0
	}

	return nil
}
