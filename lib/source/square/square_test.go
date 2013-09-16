package square

import (
	"github.com/bigwill/process/core"
	"testing"
)

func TestIsSource(t *testing.T) {
	ctx := core.NewContext(48000, 2, 100)
	s, err := NewSource(ctx)

	if err != nil {
		t.Fatal(err)
	}

	switch s.(type) {
	case core.Source:
		t.Log("type is correct")
	default:
		t.Fatal("package does not implement Source interface")
	}
}
