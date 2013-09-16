package stdout

import (
	"github.com/bigwill/process/core"
	"testing"
)

func TestIsSink(t *testing.T) {
	ctx := core.NewContext(48000, 2, 100)
	s, err := NewSink(ctx)

	if err != nil {
		t.Fatal(err)
	}

	switch s.(type) {
	case core.Sink:
		t.Log("type is correct")
	default:
		t.Fatal("package does not implement Sink interface")
	}
}
