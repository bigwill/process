package filter

import (
	"github.com/bigwill/process/core"
	"github.com/bigwill/process/lib/params/linear"
	"math"
)

const (
	LPF12 = iota
	HPF12 = iota
	BPF6  = iota
	NF2P  = iota
	NumFilterTypes
)

type State struct {
	sampleRate core.Quantity
	f_type     core.Param
	f_c        core.Param
	f_Q        core.Param

	// Parameter Values
	filterType int
	cutoff     core.Quantity
	q          core.Quantity

	// sample state
	x_n1 core.Quantity
	x_n2 core.Quantity
	y_n1 core.Quantity
	y_n2 core.Quantity

	// filter coefficients
	b0Norm core.Quantity
	b1Norm core.Quantity
	b2Norm core.Quantity
	a1Norm core.Quantity
	a2Norm core.Quantity
}

func NewProcessor(sampleRate core.Quantity) core.Processor {
	s := &State{sampleRate: sampleRate,
		f_type: linear.NewState("Mode", "", 0, NumFilterTypes-.01, 0.0),
		f_c:    linear.NewState("Cutoff", "Hz", 30.0, 20000.0, .1),
		f_Q:    linear.NewState("Q", "", .1, 18, .2)}

	s.f_type.SetHandler(func(p core.Param) {
		s.setType(p.Val())
	})
	s.f_type.SetPos(0)

	s.f_c.SetHandler(func(p core.Param) {
		s.setCutoff(p.Val())
	})
	s.f_c.SetPos(.5)

	s.f_Q.SetHandler(func(p core.Param) {
		s.setQ(p.Val())
	})
	s.f_Q.SetPos(.1)

	return s
}

func (s *State) Name() string {
	return "Filter"
}

func (s *State) NumParams() core.ParamIdx {
	return 3
}

func (s *State) Param(idx core.ParamIdx) core.Param {
	switch idx {
	case 0:
		return s.f_type
	case 1:
		return s.f_c
	case 2:
		return s.f_Q
	default:
		return nil
	}
}

func (s *State) Process(x_n core.Quantity) (core.Quantity, error) {
	// Processing
	y_n := x_n*s.b0Norm + s.x_n1*s.b1Norm + s.x_n2*s.b2Norm - s.y_n1*s.a1Norm - s.y_n2*s.a2Norm
	// Update state for next sample
	s.y_n2 = s.y_n1
	s.y_n1 = y_n
	s.x_n2 = s.x_n1
	s.x_n1 = x_n

	return y_n, nil
}

func (s *State) setType(t core.Quantity) {
	s.filterType = int(t)
	s.computeCoefficients()
}

func (s *State) setCutoff(cutoff core.Quantity) {
	s.cutoff = cutoff
	s.computeCoefficients()
}

func (s *State) setQ(Q core.Quantity) {
	s.q = Q
	s.computeCoefficients()
}

func (s *State) computeCoefficients() {
	w0 := float64(2.0 * math.Pi * s.cutoff / s.sampleRate)
	alph := core.Quantity(math.Sin(w0)) / 2.0 / s.q

	var b0, b1, b2, a0, a1, a2 core.Quantity

	switch s.filterType {
	case LPF12:
		b1 = 1 - core.Quantity(math.Cos(w0))
		b0 = b1 / 2
		b2 = b0
		a0 = 1 + alph
		a1 = -2 * core.Quantity(math.Cos(w0))
		a2 = 1 - alph

	case HPF12:
		b1 = -1 - core.Quantity(math.Cos(w0))
		b0 = (1 + core.Quantity(math.Cos(w0))) / 2
		b2 = b0
		a0 = 1 + alph
		a1 = -2 * core.Quantity(math.Cos(w0))
		a2 = 1 - alph

	case BPF6:
		b1 = 0
		b0 = s.q * alph
		b2 = -1 * b0
		a0 = 1 + alph
		a1 = -2 * core.Quantity(math.Cos(w0))
		a2 = 1 - alph

	case NF2P:
		b0 = 1
		b1 = -2 * core.Quantity(math.Cos(w0))
		b2 = 1
		a0 = 1 + alph
		a1 = b1
		a2 = 1 - alph
	}

	s.b0Norm = b0 / a0
	s.b1Norm = b1 / a0
	s.b2Norm = b2 / a0
	s.a1Norm = a1 / a0
	s.a2Norm = a2 / a0
}
