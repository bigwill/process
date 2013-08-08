package core

type Quantity float64
type Index int64
type ParamIdx int16

type Message interface {
	Code() int16
}

type ControlMessage Message

type MonitorMessage interface {
	Message
	ParamerName() string
}

type MidiControlMessage interface {
	ControlMessage
	Midi() MidiMessage
}

type ErrorMonitorMessage interface {
	MonitorMessage
	Err() error
}

type ControlChannel chan ControlMessage
type MonitorChannel chan MonitorMessage
type SampleChannel chan Quantity

// Message codes
const (
	Quit = 0
	Error
	Midi
)

type Param interface {
	Name() string
	SetHandler(func(Param))
	// Pos must be 0.0-1.0 inclusive linear/"physical" control position
	SetPos(Quantity)
	Pos() Quantity
	// useful value in processing (e.g., 5000 for 5kHZ in filter param)
	Val() Quantity
	// user-displayable string representation
	Repr() string
}

type Paramer interface {
	Name() string
	NumParams() ParamIdx
	Param(ParamIdx) Param
}

type Source interface {
	Paramer
	Output() (Quantity, error)
}

type SourceRoutine func(out SampleChannel, ctrl ControlChannel, mon MonitorChannel)

type Processor interface {
	Paramer
	Process(Quantity) (Quantity, error)
}

type ProcessorRoutine func(in SampleChannel, out SampleChannel, ctrl ControlChannel, mon MonitorChannel)

type Sink interface {
	Paramer
	Input(Quantity) error
}

type SinkRoutine func(out SampleChannel, ctrl ControlChannel, mon MonitorChannel)
