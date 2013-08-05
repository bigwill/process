package square

import (
	"github.com/bigwill/process/core"
)

type State struct {
	T int64 // wave period in terms of samples @ the given sample rate
	i int64 // current index in wave period
}

func MakeState(sample_rate core.Quantity) *State {
	var f_g core.Quantity = 3840
	return &State{T: int64(sample_rate / f_g)}
}

func (s *State) Generate() core.Quantity {
	defer func (t *State) {
		t.i = (t.i + 1) % t.T
	}(s)
	if s.i <= s.T / 2 {
		return 1.0
	} else {
		return -1.0
	}
}