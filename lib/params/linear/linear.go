package linear

import (
	"fmt"
        "github.com/bigwill/process/core"
)

type State struct {
	name string
	units string
	handler func(core.Param)
	curPos core.Quantity
	minVal core.Quantity
	maxVal core.Quantity
	curVal core.Quantity
}

func MakeState(name string, units string, minVal core.Quantity, maxVal core.Quantity, curPos core.Quantity) *State {
	s := &State{name: name, units: units, minVal: minVal, maxVal: maxVal, curPos: curPos}
	s.SetPos(curPos)
	return s
}

func (s *State) Name() string {
	return s.name
}

func (s *State) SetHandler(aHandler func(core.Param)) {
	s.handler = aHandler
}

func (s *State) SetPos(p core.Quantity) {
	s.curVal = core.Quantity(p) * (s.maxVal - s.minVal)
	if s.handler != nil {
		s.handler(s)
	}
}

func (s *State) Pos() core.Quantity {
	return s.curPos
}

func (s *State) Val() core.Quantity {
	return s.curVal
}

func (s *State) Repr() string {
	return fmt.Sprintf("%f %s", s.curVal, s.units)
}
