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
	"os"
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
	nnodes, err := GetNumasList()
	if err != nil {
		return 0, err
	}
	return len(nnodes), nil
}

// GetNumasList return list of numa nodes names
func GetNumasList() ([]string, error) {
	f, err := os.Open(sysDevicesSystemNodePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", sysDevicesSystemNodePath, err)
	}
	files, err := f.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("failed to list items under %s: %v", sysDevicesSystemNodePath, err)
	}

	re, err := regexp.Compile("node[0-9]*")
	if err != nil {
		return nil, err
	}

	nnodes := []string{}
	for _, f := range files {
		match := re.MatchString(f.Name())
		if match {
			nnodes = append(nnodes, f.Name())
		}
	}
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
		out, err := os.ReadFile(cpuListPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get content of %s: %v", cpuListPath, err)
		}

		numaToCpu[idx], err = cpuset.Parse(strings.TrimSpace(string(out)))
		if err != nil {
			return nil, fmt.Errorf("could not parse numa cpuset: %v", err)
		}
	}

	return numaToCpu, nil
}

// GetNumaDeviceMapping return pci-device -> numa mapping, e.g ["0000:00:00.0":0, "0000:00:02.0":1,"0000:00:04.0":0]
func GetNumaDeviceMapping() (map[string]int, error) {
	numaCount, err := GetNumaCount()
	if err != nil {
		return nil, err
	}
	numaToPci := make(map[string]int, numaCount)

	f, err := os.Open(sysBusPciDevicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %v", sysBusPciDevicePath, err)
	}
	devicesNames, err := f.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("failed to list items under %s: %v", sysBusPciDevicePath, err)
	}

	for _, d := range devicesNames {
		dPath := filepath.Join(sysBusPciDevicePath, d.Name())
		numaPath := filepath.Join(dPath, "numa_node")
		out, err := os.ReadFile(numaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to get numa node of pci device %s: %v", d.Name(), err)
		}
		nnode, err := strconv.Atoi(string(out))
		if err != nil {
			return nil, fmt.Errorf("could not parse numa of pci device %s: %v", d.Name(), err)
		}

		if nnode == -1 {
			continue
		}

		numaToPci[d.Name()] = nnode
	}

	return numaToPci, nil
}
