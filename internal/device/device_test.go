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
		testMap      map[string]int
		devList      []string
		expectedNuma int
	}{
		{
			devNumaMap,
			[]string{"example.com/deviceA"},
			0,
		},
		{
			devNumaMap,
			[]string{"example.com/deviceA", "example.com/deviceB"},
			-1,
		},
		{
			devNumaMap,
			[]string{"example.com/deviceC", "example.com/deviceB"},
			1,
		},
	}

	for _, c := range testCases {
		out := NumaAlignmentOutput{NNode: -1}
		CheckPciDeviceToNumaMapping(c.testMap, c.devList, &out)
		if out.NNode != c.expectedNuma {
			t.Fatalf("expected: %d, actual: %d, devices list: [%v]", c.expectedNuma, out.NNode, c.devList)
		}
	}

}
