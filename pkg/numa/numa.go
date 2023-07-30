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
	"strconv"
	"strings"

	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const (
	sysDevicesSystemNodePath = "/sys/devices/system/node/"
	sysBusPciDevicePath      = "/sys/bus/pci/devices"
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
		fmt.Println(string(out[:])) //TODO debug print
		numaToCpu[idx], err = cpuset.Parse(strings.TrimSpace(string(out[:])))
		if err != nil {
			return nil, fmt.Errorf("could not parse numa cpuset: %v", err)
		}
	}

	return numaToCpu, nil
}

// GetNumaDeviceMapping return pci-device -> numa mapping, e.g node0: ["0000:00:00.0","0000:00:02.0","0000:00:04.0"]
func GetNumaDeviceMapping() (map[int][]string, error) {
	numaCount, err := GetNumaCount()
	if err != nil {
		return nil, err
	}
	numaToPci := make(map[int][]string, numaCount)

	out, err := exec.Command("ls", sysBusPciDevicePath).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list items under %s: %v", sysBusPciDevicePath, err)
	}
	devicesNames := strings.Fields(string(out[:]))
	fmt.Printf("devices: \n %v \n ", devicesNames) //TODO debug print
	for _, dName := range devicesNames {
		dPath := filepath.Join(sysBusPciDevicePath, dName)
		numaPath := filepath.Join(dPath, "numa_node")
		out, err := exec.Command("cat", numaPath).Output()
		if err != nil {
			return nil, fmt.Errorf("failed to get numa node of pci device %s: %v", dName, err)
		}
		nnode, err := strconv.Atoi(string(out[:]))
		if err != nil {
			return nil, fmt.Errorf("could not parse numa of pci device %s: %v", dName, err)
		}

		if nnode == -1 {
			continue
		}

		numaToPci[nnode] = append(numaToPci[nnode], dName)
	}

	return numaToPci, nil
}
