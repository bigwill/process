package sox

import (
	"encoding/binary"
	"github.com/bigwill/process/core"
	"io"
	"log"
	"os/exec"
)

const bufferSize = 500

type State struct {
	i int
	buf []core.Quantity
	reader io.Reader
}

func NewSource(sampleRate core.Quantity) core.Source {
	s := &State{buf: make([]core.Quantity, bufferSize, bufferSize)}

	// TODO: yuck. generalize this at some point
	cmd := exec.Command("sox", "~/Downloads/VST/Crumar_Cello.wav", "-t", "f64", "-r", "48k", "-c", "1", "-")

	var err error
	s.reader, err = cmd.StdoutPipe()
	if err != nil { // TODO: better error reporting
		log.Printf("cmd fail = %v", err)
		return nil
	}

	err = cmd.Start()
	if err != nil { // TODO: better error reporting
		log.Printf("cmd fail = %v", err)
		return nil
	}

	return s
}

func (s *State) Name() string {
	return "Sox"
}

func (s *State) NumParams() core.ParamIdx {
	return 0
}

func (s *State) Param(idx core.ParamIdx) core.Param {
	return nil
}

func (s *State) Output() (core.Quantity, error) {
	if s.i == len(s.buf) {
		err := binary.Read(s.reader, binary.LittleEndian, s.buf)
		if err != nil {
			return 0, err
		}

		s.i = 0
	}

	defer func() { s.i++ }()
	return s.buf[s.i], nil
}
