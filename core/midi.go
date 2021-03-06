package core

import (
	"encoding/binary"
	"io"
)

const (
	NoteOff          = 128
	NoteOn           = 144
	Aftertouch       = 160
	ControlChange    = 176
	ProgramChange    = 192
	ChannelPressure  = 208
	PitchWheelChange = 224
)

type MidiMessage interface {
	Event() byte
	Channel() byte
	Key() byte
	Value() byte
	Program() byte
	Pitch() uint16
}

type state struct {
	event   byte
	channel byte
	key     byte
	value   byte
	program byte
	pitch   uint16
}

func (s state) Event() byte {
	return s.event
}

func (s state) Channel() byte {
	return s.channel
}

func (s state) Key() byte {
	return s.key
}

func (s state) Value() byte {
	return s.value
}

func (s state) Program() byte {
	return s.program
}

func (s state) Pitch() uint16 {
	return s.pitch
}

func newMidiMessage(e byte, ch byte, key byte, val byte, prog byte, pitch uint16) MidiMessage {
	return state{event: e, channel: ch, key: key, value: val, program: prog, pitch: pitch}
}

func ReadMidi(r io.Reader) (MidiMessage, error) {
	for {
		var evChan, event, channel, key, value uint8

		err := binary.Read(r, binary.LittleEndian, &evChan)
		if err != nil {
			return nil, err
		}

		channel = evChan & 0x0f
		event = evChan & 0xf0

		switch event {
		case NoteOff:
			fallthrough
		case NoteOn:
			fallthrough
		case Aftertouch:
			fallthrough
		case ControlChange:
			err := binary.Read(r, binary.LittleEndian, &key)
			if err != nil {
				return nil, err
			}

			err = binary.Read(r, binary.LittleEndian, &value)
			if err != nil {
				return nil, err
			}

			return newMidiMessage(event, channel, key, value, 0, 0), nil
		case ProgramChange: // TODO: handle program change
			fallthrough
		case ChannelPressure: // TODO: handle channel pressure
			var scratch byte
			err := binary.Read(r, binary.LittleEndian, &scratch)
			if err != nil {
				return nil, err
			}

		case PitchWheelChange: // TODO: handle pitch wheel
			var scratch uint16
			err := binary.Read(r, binary.LittleEndian, &scratch)
			if err != nil {
				return nil, err
			}
		}
	}
}
