package tests

import (
	"testing"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	"github.com/shajmakh/numaalign-rewritten/internal/cpu"
	"github.com/shajmakh/numaalign-rewritten/internal/device"
	"github.com/shajmakh/numaalign-rewritten/internal/memory"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

func TestResourcesNumaAlign(t *testing.T) {
	cpuMap := map[int]cpuset.CPUSet{
		0: cpu.GetCpuset("0-2,6-7,16-23"),
		1: cpu.GetCpuset("3-5,8-15"),
	}

	pciMap := map[string]int{
		"devA": 1,
		"devB": 1,
		"devC": 0,
	}

	testCases := []struct {
		description       string
		consumedResources ProccessResources
		numaToCpuMap      map[int]cpuset.CPUSet
		pciToNumaMap      map[string]int
		expectedNuma      int
		expectedIsAligned bool
	}{
		{
			description: "aligned",
			consumedResources: ProccessResources{
				CPUs:   cpu.GetCpuset("3-5,11,13"),
				PCI:    []string{"devA", "devB"},
				Memory: "1",
			},
			numaToCpuMap:      cpuMap,
			pciToNumaMap:      pciMap,
			expectedNuma:      1,
			expectedIsAligned: true,
		},
		{
			description: "aligned",
			consumedResources: ProccessResources{
				CPUs:   cpu.GetCpuset("0,2"),
				PCI:    []string{},
				Memory: "0",
			},
			numaToCpuMap:      cpuMap,
			pciToNumaMap:      pciMap,
			expectedNuma:      0,
			expectedIsAligned: true,
		},
		{
			description: "not aligned (memory)",
			consumedResources: ProccessResources{
				CPUs:   cpu.GetCpuset("0,2"),
				PCI:    []string{},
				Memory: "0-1",
			},
			numaToCpuMap:      cpuMap,
			pciToNumaMap:      pciMap,
			expectedNuma:      -1,
			expectedIsAligned: false,
		},
		{
			description: "not aligned (devices)",
			consumedResources: ProccessResources{
				CPUs:   cpu.GetCpuset("3-5"),
				PCI:    []string{"devA", "devC"},
				Memory: "1",
			},
			numaToCpuMap:      cpuMap,
			pciToNumaMap:      pciMap,
			expectedNuma:      -1,
			expectedIsAligned: false,
		},
		{
			description: "not aligned (cpus)",
			consumedResources: ProccessResources{
				CPUs:   cpu.GetCpuset("1-6"),
				PCI:    []string{"devC"},
				Memory: "0",
			},
			numaToCpuMap:      cpuMap,
			pciToNumaMap:      pciMap,
			expectedNuma:      -1,
			expectedIsAligned: false,
		},
		{
			description: "not aligned (memory)",
			consumedResources: ProccessResources{
				CPUs:   cpu.GetCpuset("1"),
				PCI:    []string{"devC"},
				Memory: "1",
			},
			numaToCpuMap:      cpuMap,
			pciToNumaMap:      pciMap,
			expectedNuma:      -1,
			expectedIsAligned: false,
		},
		{
			description: "not aligned (memory)",
			consumedResources: ProccessResources{
				CPUs:   cpu.GetCpuset("1"),
				PCI:    []string{"devA", "devB"},
				Memory: "0-1",
			},
			numaToCpuMap:      cpuMap,
			pciToNumaMap:      pciMap,
			expectedNuma:      -1,
			expectedIsAligned: false,
		},
	}

	for _, tc := range testCases {
		numa, isAligned := cpuPciIntegrationAlignment(tc.consumedResources, tc.numaToCpuMap, tc.pciToNumaMap)
		if isAligned != tc.expectedIsAligned || numa != tc.expectedNuma {
			t.Fatalf("expected alignment: %t:%d ; actual: %t:%d ; cpuset: [%s], devices list: [%v], memory nodes: [%s]", tc.expectedIsAligned, tc.expectedNuma, isAligned, numa, tc.consumedResources.CPUs, tc.consumedResources.PCI, tc.consumedResources.Memory)
		}
	}
}

func cpuPciIntegrationAlignment(res ProccessResources, cpuMap map[int]cpuset.CPUSet, pciMap map[string]int) (int, bool) {
	out := NewOutput()
	// till now we do not have a cheap way to write e2e real test cases, thus we simulate the steps
	// that main does to check the alignments, taking into account the order of the checks.
	// Hence, these tests does not provide full tests coverage.
	cpu.CheckNumaCpuMapping(cpuMap, res.CPUs, &out)
	device.CheckPciDeviceToNumaMapping(pciMap, res.PCI, &out)
	memory.CheckAlignmentWith(res.Memory, &out)
	return out.NNode, out.IsAligned
}
