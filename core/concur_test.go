package core

import (
	"testing"
)

const (
	concurSampleRate = 48000
	concurChannels = 2
	concurPoolSize = 100
)

//
// Test source
//
type srcState struct {
	i   int
	spyChan SampleChannel
}

func newTestSource() *srcState {
	return &srcState{spyChan:make(SampleChannel)}
}

func (s *srcState) Name() string {
	return "Test Source"
}

func (s *srcState) NumParams() ParamIdx {
	return 0
}

func (s *srcState) Param(idx ParamIdx) Param {
	return nil
}

func (s *srcState) Output(fr SampleFrame) error {
	var v Quantity
	if s.i == 0 {
		v = 1
	} else {
		v = -1
	}

	for j := Integer(0); j < fr.NumChannels(); j++ {
		fr.SetChannelVal(j, v)
		fr.SetChannelVal(j, -1*v)
	}

	s.i++
	s.spyChan <- fr

	return nil
}

func (s *srcState) spyChannel() SampleChannel {
	return s.spyChan
}

//
// Test processor
//
type procState struct {
}

func newTestProcessor() Processor {
	return &procState{}
}

func (s *procState) Name() string {
	return "Test Proc"
}

func (s *procState) NumParams() ParamIdx {
	return 0
}

func (s *procState) Param(idx ParamIdx) Param {
	return nil
}

func (s *procState) Process(x SampleFrame) error {
	// leave the sample be
	return nil
}

//
// Test sink
//
type snkState struct {
	i   int
	spyChan SampleChannel
}

func newTestSink() *snkState {
	return &snkState{spyChan:make(SampleChannel)}
}

func (s *snkState) Name() string {
	return "Test Sink"
}

func (s *snkState) NumParams() ParamIdx {
	return 0
}

func (s *snkState) Param(idx ParamIdx) Param {
	return nil
}

func (s *snkState) Input(fr SampleFrame) error {
	s.spyChan <- fr
	return nil
}

func (s *snkState) spyChannel() SampleChannel {
	return s.spyChan
}

func TestRunChainNoMidi(t *testing.T) {
	ctx := NewContext(concurSampleRate, concurChannels, concurPoolSize)
	src := newTestSource()
	procs := []Processor{newTestProcessor(), newTestProcessor(), newTestProcessor()}
	snk := newTestSink()

	RunChain(ctx, src, procs, snk, nil)

	for i := 0; i < 100; i++ {
		f1 := <- src.spyChannel()
		f2 := <- snk.spyChannel()
		for j := Integer(0); j < ctx.NumChannels(); j++ {
			if f1.ChannelVal(j) == f2.ChannelVal(j) {
				t.Logf("%v'th channel of %v'th frame matched (src=%v, snk=%v)", j, i, f1.ChannelVal(j), f2.ChannelVal(j))
			} else {
				t.Fatalf("%v'th channel of %v'th frame didn't match as expected (src=%v, snk=%v)", j, i, f1.ChannelVal(j), f2.ChannelVal(j))
			}
		}
	}
}