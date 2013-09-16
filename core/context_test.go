package core

import (
	"testing"
)

const (
	sampleRate = 48000
	numChannels = 2
	framePoolSize = 100
)

func TestContextGetters (t *testing.T) {
	ctx := NewContext(sampleRate, numChannels, framePoolSize)

	gSampleRate := ctx.SampleRate()
	if gSampleRate != sampleRate {
		t.Error("got %v sample rate, expected %v", gSampleRate, sampleRate)
	}

	gNumChannels := ctx.NumChannels()
	if gNumChannels != numChannels {
		t.Error("got %v channels, expected %v", gNumChannels, numChannels)
	}

	gFramePoolSize := ctx.FramePool().Size()
	if gFramePoolSize != framePoolSize {
		t.Error("got %v frame pool size, expected %v", gFramePoolSize, framePoolSize)
	}
}