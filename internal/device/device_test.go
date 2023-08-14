package device

import (
	"testing"

	. "github.com/shajmakh/numaalign-rewritten/internal"
)

func TestCheckPciDevicesAlignment(t *testing.T) {

	devNumaMap := map[string]int{
		"example.com/deviceA": 0,
		"example.com/deviceB": 1,
		"example.com/deviceC": 1,
	}

	testCases := []struct {
		testMap           map[string]int
		devList           []string
		expectedNuma      int
		expectedIsAligned bool
	}{
		{
			devNumaMap,
			[]string{"example.com/deviceA"},
			0,
			true,
		},
		{
			devNumaMap,
			[]string{"example.com/deviceA", "example.com/deviceB"},
			-1,
			false,
		},
		{
			devNumaMap,
			[]string{"example.com/deviceC", "example.com/deviceB"},
			1,
			true,
		},
	}

	for _, c := range testCases {
		out := NewOutput()
		CheckPciDeviceToNumaMapping(c.testMap, c.devList, &out)
		if out.NNode != c.expectedNuma || out.IsAligned != c.expectedIsAligned {
			t.Fatalf("expected alignment: %t:%d ; actual: %t/%d ; devices list: [%v]", c.expectedIsAligned, c.expectedNuma, out.IsAligned, out.NNode, c.devList)
		}
	}

}
