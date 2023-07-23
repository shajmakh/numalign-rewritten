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

func CheckCpuAlignment(output *NumaAlignmentOutput) {

	//get numa->cpuset map
	numaToCpuset, err := GetNumaCpuMapping()
	if err != nil {
		output.Err = err
		return
	}
	//get consumed cpus
	consumedCpuset, err := GetConsumedCpusBy() //TODO send pid
	if err != nil {
		output.Err = err
		return
	}
	output.ProccessResources.CPUs = consumedCpuset

	// check mapping
	CheckNumaCpuMapping(numaToCpuset, consumedCpuset, *output)
}

func CheckNumaCpuMapping(numaToCpuset map[int]cpuset.CPUSet, consumedCpuset cpuset.CPUSet, output NumaAlignmentOutput) int {
	for idx, allocatedCpuset := range numaToCpuset {
		if consumedCpuset.IsSubsetOf(allocatedCpuset) {
			if output.NNode != -1 && output.NNode != idx {
				output.IsAligned = false
				output.NNode = -1
				break
			}
			output.IsAligned = true //the process may consume only this resource type, so it is aligned hence true
			output.NNode = idx
			break
		}
	}
	return output.NNode
}
