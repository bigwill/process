package stdout

import (
	"encoding/binary"
	"github.com/bigwill/process/core"
	"os"
)

const bufferSize = 500

type State struct {
	i   int
	buf []core.Quantity
}

func NewSink(sampleRate core.Quantity) core.Sink {
	return &State{buf: make([]core.Quantity, bufferSize, bufferSize)}
}

func (s *State) Name() string {
	return "Std Out"
}

func (s *State) NumParams() core.ParamIdx {
	return 0
}

func (s *State) Param(idx core.ParamIdx) core.Param {
	return nil
}

func (s *State) Input(v core.Quantity) error {
	s.buf[s.i] = v
	s.i++

	if s.i == bufferSize {
		err := binary.Write(os.Stdout, binary.LittleEndian, s.buf)
		if err != nil {
			return err
		}

		s.i = 0
	}

	return nil
}
