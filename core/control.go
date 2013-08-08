package core

type ctrlMsg struct {
	code int16
}

func (m *ctrlMsg) Code() int16 {
	return m.code
}

type midiCtrlMsg struct {
	ctrlMsg
	midi MidiMessage
}

func (m *midiCtrlMsg) Midi() MidiMessage {
	return m.midi
}

func NewQuitControlMessage() ControlMessage {
	return &ctrlMsg{Quit}
}

func NewMidiControlMessage(m MidiMessage) ControlMessage {
	return &midiCtrlMsg{ctrlMsg{Midi}, m}
}