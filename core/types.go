package core

type Quantity float64
type Index int64
type ParamIdx int16

type Message interface {
	Code() int16
}

type ControlMessage Message

type ProcessorContext interface {
	SampleRate() Quantity
	NumChannels() Index
	FramePool() SampleFramePool
}

type MonitorMessage interface {
	Message
	ParamerName() string
}

type ParamControlMessage interface {
	ControlMessage
	Index() ParamIdx
	Pos() Quantity
}

type ErrorMonitorMessage interface {
	MonitorMessage
	Err() error
}

type ControlChannel chan ControlMessage
type MonitorChannel chan MonitorMessage
type SampleChannel chan SampleFrame

type SampleFrame interface {
	// NOT channel in golang parlance, but in audio signal parlance (e.g., stereo has 2 audio channels)
	NumChannels() Index
	ChannelVal(Index) Quantity
	SetChannelVal(Index, Quantity)
}

type SampleFramePool interface {
	DequeueFrame() SampleFrame
	EnqueueFrame(SampleFrame)
}

// Message codes
const (
	Quit        = iota
	Error       = iota
	ParamChange = iota
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
	Output() (SampleFrame, error)
}

type Processor interface {
	Paramer
	Process(SampleFrame) (SampleFrame, error)
}

type Sink interface {
	Paramer
	Input(SampleFrame) error
}

type MidiSource interface {
	Name() string
	Output() (MidiMessage, error)
	Mapper() MidiMapper
}

type MidiMapper interface {
	Map(MidiMessage) ParamControlMessage
}
