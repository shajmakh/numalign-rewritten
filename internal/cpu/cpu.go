/*
 *Copyright 2023 Red Hat, Inc.
 *
 *Licensed under the Apache License, Version 2.0 (the "License");
 *you may not use this file except in compliance with the License.
 *You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 *Unless required by applicable law or agreed to in writing, software
 *distributed under the License is distributed on an "AS IS" BASIS,
 *WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *See the License for the specific language governing permissions and
 *limitations under the License.
 */

package cpu

import (
	"fmt"
	"os/exec"
	"strings"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	. "github.com/shajmakh/numaalign-rewritten/pkg/numa"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

// GetConsumedCpusBy returns the consumed cpuset by a proccess
func GetConsumedCpusBy(pid string) (cpuset.CPUSet, error) {
	var consumedCpuset cpuset.CPUSet
	out, err := exec.Command("grep", "Cpus_allowed_list", fmt.Sprintf("/proc/%s/status", pid)).Output()
	if err != nil {
		return consumedCpuset, fmt.Errorf("could not get the status of process %s: %v", pid, err)
	}

	val := strings.Split(string(out[:]), ":")[1]
	consumedCpuset, err = cpuset.Parse(strings.TrimSpace(val))
	if err != nil {
		return consumedCpuset, fmt.Errorf("could not parse cpuset: %v", err)
	}

	return consumedCpuset, nil
}

// CheckCpuAlignment checks if cpus consumed by a process are aligned to a single numa node
func CheckCpuAlignment(pid string, output *NumaAlignmentOutput) {
	numaToCpuset, err := GetNumaCpuMapping()
	if err != nil {
		output.Err = err
		return
	}

	consumedCpuset, err := GetConsumedCpusBy(pid)
	if err != nil {
		output.Err = err
		return
	}
	output.ProccessResources.CPUs = consumedCpuset

	CheckNumaCpuMapping(numaToCpuset, consumedCpuset, output)
}

// CheckNumaCpuMapping checks if a cpuset is mapped to a numa and returns that numa
func CheckNumaCpuMapping(numaToCpuset map[int]cpuset.CPUSet, consumedCpuset cpuset.CPUSet, output *NumaAlignmentOutput) {
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
}
