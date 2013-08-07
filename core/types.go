package core

type Quantity float64
type Index int64
type ParamIdx int16
type Control int16

// Control constants
const (
	Quit = 0
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

type Generator interface {
	Paramer
	Generate() Quantity
}

type Processor interface {
	Paramer
	Process(Quantity) Quantity
}
