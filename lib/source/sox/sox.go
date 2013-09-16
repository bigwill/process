package sox

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
	reader io.Reader
}

func NewSource(ctx core.Context) (core.Source, error) {
	s := &state{ctx: ctx, buf: make([]core.Quantity, bufferSize*ctx.NumChannels())}

	// TODO: yuck. generalize this at some point
	cmd := exec.Command("sox", "~/go/process/src/github.com/bigwill/process/res/Crumar_Cello.wav", "-t", "f64", "-r", fmt.Sprintf("%v", ctx.SampleRate()), "-c", fmt.Sprintf("%v", ctx.NumChannels()), "-")

	var err error
	s.reader, err = cmd.StdoutPipe()
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
	return "Sox"
}

func (s *state) NumParams() core.ParamIdx {
	return 0
}

func (s *state) Param(idx core.ParamIdx) core.Param {
	return nil
}

func (s *state) Output(fr core.SampleFrame) error {
	if s.i == len(s.buf) {
		err := binary.Read(s.reader, binary.LittleEndian, s.buf)
		if err != nil {
			return err
		}

		s.i = 0
	}

	for j := core.Index(0); j < s.ctx.NumChannels(); j++ {
		fr.SetChannelVal(j, s.buf[s.i])
		s.i++
	}

	return nil
}
