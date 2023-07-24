package cpu

import (
	"fmt"
	"os/exec"
	"strings"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	. "github.com/shajmakh/numaalign-rewritten/pkg/numa"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const (
	sysFsCgroupCpusetCpusPath = "/sys/fs/cgroup/cpuset/cpuset.cpus"
)

// GetConsumedCpusBy returns the consumed cpuset by a proccess
func GetConsumedCpusBy() (cpuset.CPUSet, error) {
	var consumedCpuset cpuset.CPUSet
	out, err := exec.Command("cat", sysFsCgroupCpusetCpusPath).Output()
	if err != nil {
		return consumedCpuset, fmt.Errorf("could not list cpuset from %s: %v", sysFsCgroupCpusetCpusPath, err)
	}

	consumedCpuset, err = cpuset.Parse(strings.TrimSpace(string(out[:])))
	if err != nil {
		return consumedCpuset, fmt.Errorf("could not parse cpuset: %v", err)
	}

	return consumedCpuset, nil
}

// CheckCpuAlignment checks if cpus consumed by a process are aligned to a single numa node
func CheckCpuAlignment(output *NumaAlignmentOutput) {
	numaToCpuset, err := GetNumaCpuMapping()
	if err != nil {
		output.Err = err
		return
	}

	consumedCpuset, err := GetConsumedCpusBy() //TODO send pid
	if err != nil {
		output.Err = err
		return
	}
	output.ProccessResources.CPUs = consumedCpuset

	CheckNumaCpuMapping(numaToCpuset, consumedCpuset, *output)
}

// CheckNumaCpuMapping checks if a cpuset is mapped to a numa and returns that numa
func CheckNumaCpuMapping(numaToCpuset map[int]cpuset.CPUSet, consumedCpuset cpuset.CPUSet, output NumaAlignmentOutput) int {
	for idx, allocatedCpuset := range numaToCpuset {
		if consumedCpuset.IsSubsetOf(allocatedCpuset) {
			if output.NNode != -1 && output.NNode != idx {
				output.NNode = -1
				break
			}
			output.NNode = idx
			break
		}
	}
	return output.NNode
}
