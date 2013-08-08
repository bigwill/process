package play

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
	writer io.Writer
}

func NewSink(sampleRate core.Quantity) core.Sink {
	s := &State{buf: make([]core.Quantity, bufferSize, bufferSize)}

	// TODO: yuck. make this cleaner
	cmd := exec.Command("play", "-t", "f64", "-r", "48k", "-c", "1", "-")

	var err error
	s.writer, err = cmd.StdinPipe()
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
	return "Sox Play"
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

	if (s.i == bufferSize) {
		err := binary.Write(s.writer, binary.LittleEndian, s.buf)
		if err != nil {
			return err
		}

		s.i = 0
	}

	return nil
}
