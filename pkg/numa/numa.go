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

package numa

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const (
	sysDevicesSystemNodePath = "/sys/devices/system/node/"
)

// GetNumaCount returns the numa node's count on the system
func GetNumaCount() (int, error) {
	//out, err := exec.Command("lscpu", "--json  |jq '.[] '").Output() //| .[] | select(.field==\"NUMA node(s):\")| .data'").Output()
	nnodes, err := GetNumasList()
	if err != nil {
		return 0, err
	}
	return len(nnodes), nil
}

// GetNumasList return list of numa nodes names
func GetNumasList() ([]string, error) {
	out, err := exec.Command("ls", sysDevicesSystemNodePath).Output()
	if err != nil {
		return []string{}, fmt.Errorf("failed to list items under %s: %v", sysDevicesSystemNodePath, err)
	}

	nNodeDirRegex := regexp.MustCompile("node[0-9]*")
	nnodes := nNodeDirRegex.FindAllString(string(out[:]), -1)

	return nnodes, nil
}

// GetNumaCpuMapping return cpuset -> numa mapping, e.g node0: {0,5-8,12}
func GetNumaCpuMapping() (map[int]cpuset.CPUSet, error) {
	numaToCpu := make(map[int]cpuset.CPUSet)

	nnodes, err := GetNumasList()
	if err != nil {
		return nil, err
	}

	for idx, nnode := range nnodes {
		nnodePath := filepath.Join(sysDevicesSystemNodePath, nnode)
		cpuListPath := filepath.Join(nnodePath, "cpulist")
		out, err := exec.Command("cat", cpuListPath).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to list items under %s: %v", cpuListPath, err)
		}
		fmt.Println(string(out[:]))
		numaToCpu[idx], err = cpuset.Parse(strings.TrimSpace(string(out[:])))
		if err != nil {
			return nil, fmt.Errorf("could not parse numa cpuset: %v", err)
		}
	}

	return numaToCpu, nil
}
