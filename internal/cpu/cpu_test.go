package cpu

import (
	"testing"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func TestCheckNumaCpuMapping(t *testing.T) {

	numaCpuMap := map[int]cpuset.CPUSet{
		0: getCpuset("3-5,8-15"),
		1: getCpuset("0-2"),
	}

	testCases := []struct {
		testMap      map[int]cpuset.CPUSet
		cpuset       cpuset.CPUSet
		expectedNuma int
	}{
		{numaCpuMap, getCpuset("0-2"), 1},
		{numaCpuMap, getCpuset("5"), 0},
		{numaCpuMap, getCpuset("1,5,9,12"), -1},
		{numaCpuMap, getCpuset("0-2,1,5"), -1},
	}

	for _, c := range testCases {
		actual := CheckNumaCpuMapping(c.testMap, c.cpuset, NumaAlignmentOutput{NNode: -1})
		if actual != c.expectedNuma {
			t.Fatalf("expected: %d, actual: %d, cpuset: [%v]", c.expectedNuma, actual, c.cpuset)
		}
	}
}

func getCpuset(set string) cpuset.CPUSet {
	cpuset, _ := cpuset.Parse(set)
	return cpuset
}
