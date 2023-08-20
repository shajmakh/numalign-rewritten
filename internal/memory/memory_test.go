package memory

import (
	"testing"

	. "github.com/shajmakh/numaalign-rewritten/internal"
)

func TestCheckAlignmentWith(t *testing.T) {
	testCases := []struct {
		memoryNodeString  string
		expectedNuma      int
		expectedIsAligned bool
	}{
		{
			"0",
			0,
			true,
		},
		{
			"   0-3,5 ",
			-1,
			false,
		},
		{
			"0-1",
			-1,
			false,
		},
	}

	for _, c := range testCases {
		out := NewOutput()
		CheckAlignmentWith(c.memoryNodeString, &out)
		if out.NNode != c.expectedNuma || out.IsAligned != c.expectedIsAligned {
			t.Fatalf("expected alignment: %t:%d ; actual: %t/%d ; memory nodes:[%s]", c.expectedIsAligned, c.expectedNuma, out.IsAligned, out.NNode, c.memoryNodeString)
		}
	}

}
