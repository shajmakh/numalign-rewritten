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
	"os"
	"strings"

	"github.com/shajmakh/numaalign-rewritten/internal"
	"github.com/shajmakh/numaalign-rewritten/pkg/numa"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const CPUS_ALLOWED_LIST = "Cpus_allowed_list"

// GetConsumedCpusBy returns the consumed cpuset by a proccess
func GetConsumedCpusBy(pid string) (cpuset.CPUSet, error) {
	var consumedCpuset cpuset.CPUSet
	out, err := os.ReadFile(fmt.Sprintf("/proc/%s/status", pid))
	if err != nil {
		return consumedCpuset, fmt.Errorf("could not get the status of process %s: %v", pid, err)
	}
	outStr := string(out)

	//compile regex and get the matching output
	match := internal.GetValue(CPUS_ALLOWED_LIST, outStr)
	if len(match) == 0 {
		return consumedCpuset, fmt.Errorf("value %s not found in %s", CPUS_ALLOWED_LIST, outStr)
	}
	consumedCpuset, err = cpuset.Parse(strings.TrimSpace(match[1]))
	if err != nil {
		return consumedCpuset, fmt.Errorf("could not parse cpuset: %v", err)
	}
	return consumedCpuset, nil
}

// CheckCpuAlignment checks if cpus consumed by a process are aligned to a single numa node
func CheckCpuAlignment(pid string, output *internal.NumaAlignmentOutput) {
	numaToCpuset, err := numa.GetNumaCpuMapping()
	if err != nil {
		output.IsAligned = false
		output.Err = err
		return
	}

	consumedCpuset, err := GetConsumedCpusBy(pid)
	if err != nil {
		output.IsAligned = false
		output.Err = err
		return
	}
	output.ProccessResources.CPUs = consumedCpuset

	CheckNumaCpuMapping(numaToCpuset, consumedCpuset, output)
}

// CheckNumaCpuMapping checks if a cpuset is mapped to a numa and returns that numa
func CheckNumaCpuMapping(numaToCpuset map[int]cpuset.CPUSet, consumedCpuset cpuset.CPUSet, output *internal.NumaAlignmentOutput) {
	if !output.IsAligned {
		return
	}
	//there is no way that the consumed cpus are empty, but let's cover that case for the sake of completness
	if consumedCpuset.IsEmpty() {
		return
	}
	for idx, allocatedCpuset := range numaToCpuset {
		if consumedCpuset.IsSubsetOf(allocatedCpuset) {
			if output.NNode != -1 && output.NNode != idx {
				output.NNode = -1
				output.IsAligned = false
				return
			}
			output.NNode = idx
			return
		}
	}
	// this is a case where the consumed subset is part of more than one numa, we know it's not empty
	// so it definitely indicates more than 1 numa
	output.IsAligned = false
}

// GetCpuset converts the string "set" to CPUSet and returns it
func GetCpuset(set string) cpuset.CPUSet {
	cpuset, _ := cpuset.Parse(set)
	return cpuset
}
