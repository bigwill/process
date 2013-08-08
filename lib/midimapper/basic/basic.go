package basic

import (
	"github.com/bigwill/process/core"
)

type State struct {
}

func NewMidiMapper() core.MidiMapper {
	return &State{}
}

func (s *State) Map(midiMsg core.MidiMessage) core.ParamControlMessage {
	return core.NewParamControlMessage(core.ParamIdx(midiMsg.Key()),
		core.Quantity(midiMsg.Value())/127.0)
}