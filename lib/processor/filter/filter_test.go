package filter

import (
	"github.com/bigwill/process/core"
	"testing"
)

func TestIsProcessor(t *testing.T) {
	ctx := core.NewContext(48000, 2, 100)
	p , err := NewProcessor(ctx)

	if err != nil {
		t.Fatal(err)
	}

	switch p.(type) {
	case core.Processor:
		t.Log("type is correct")
	default:
		t.Fatal("package does not implement Processor interface")
	}
}
