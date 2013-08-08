package ps3midi

import (
	"github.com/bigwill/process/core"
	"io"
	"log"
	"os/exec"
)

type State struct {
	reader io.Reader
	mapper core.MidiMapper
}

func NewMidiSource(mapper core.MidiMapper) core.MidiSource {
	cmd := exec.Command("/Users/will/git/ps3midi/ps3midi.py", "lightsaber")

	reader, err := cmd.StdoutPipe()
	if err != nil { // TODO: better error reporting
		log.Printf("cmd fail = %v", err)
		return nil
	}

	err = cmd.Start()
	if err != nil { // TODO: better error reporting
		log.Printf("cmd fail = %v", err)
		return nil
	}

	return &State{reader, mapper}
}

func (s *State) Name() string {
	return "ps3midi"
}

func (s *State) Mapper() core.MidiMapper {
	return s.mapper
}

func (s *State) Output() (core.MidiMessage, error) {
	return core.ReadMidi(s.reader)
}
