package stdin

import (
	"encoding/binary"
	"github.com/bigwill/process/core"
	"os"
)

const bufferSize = 500

type State struct {
	i int
	buf []core.Quantity
}

func NewSource(sampleRate core.Quantity) core.Source {
	return &State{buf: make([]core.Quantity, bufferSize, bufferSize)}
}

func (s *State) Name() string {
	return "Std In"
}

func (s *State) NumParams() core.ParamIdx {
	return 0
}

func (s *State) Param(idx core.ParamIdx) core.Param {
	return nil
}

func (s *State) Output() (core.Quantity, error) {
	if s.i == len(s.buf) {
		err := binary.Read(os.Stdin, binary.LittleEndian, s.buf)
		if err != nil {
			return 0, err
		}

		s.i = 0
	}

	defer func() { s.i++ }()
	return s.buf[s.i], nil
}
