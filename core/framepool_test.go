package core

import (
	"testing"
)

const (
	size = 100
	channels = 2
)

func TestFramePool(t *testing.T) {
	fp := NewFramePool(size, channels)

	// basics
	if fp.Size() != size {
		t.Fatal("got %v pool size, expected %v", fp.Size(), size)
	}

	if fp.NumAvailable() != size {
		t.Fatal("got %v available frames, expected %v", fp.NumAvailable(), size)
	}

	// test dequeue
	f := fp.DequeueFrame()
	if fp.Size() != size {
		t.Fatal("got %v pool size, expected %v", fp.Size(), size)
	}

	if fp.NumAvailable() != size - 1 {
		t.Fatal("got %v available frames, expected %v", fp.NumAvailable(), size-1)
	}

	//test enqueue
	fp.EnqueueFrame(f)
	if fp.Size() != size {
		t.Fatal("got %v pool size, expected %v", fp.Size(), size)
	}

	if fp.NumAvailable() != size {
		t.Fatal("got %v available frames, expected %v", fp.NumAvailable(), size)
	}

	// drain
	for i := 0; i < size; i++ {
		fp.DequeueFrame()
	}

	if fp.NumAvailable() != 0 {
		t.Fatal("got %v available frames, expected %v", fp.NumAvailable(), 0)
	}
}