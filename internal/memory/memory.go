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

package memory

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	"k8s.io/kubernetes/pkg/kubelet/cm/cpuset"
)

const MEMS_ALLOWED_LIST = "Mems_allowed_list"

func CheckMemoryResourcesAlignment(pid string, output *NumaAlignmentOutput) {
	if !output.IsAligned {
		return
	}

	out, err := exec.Command("cat", fmt.Sprintf("/proc/%s/status", pid)).Output()
	if err != nil {
		output.IsAligned = false
		output.Err = err
		return
	}

	outStr := string(out[:])

	match := GetValue(MEMS_ALLOWED_LIST, outStr)
	if len(match) == 0 {
		output.IsAligned = false
		output.Err = fmt.Errorf("value %s not found in %s", MEMS_ALLOWED_LIST, outStr)
		return
	}
	CheckAlignmentWith(match[1], output)
}

func CheckAlignmentWith(memStr string, output *NumaAlignmentOutput) {
	if !output.IsAligned {
		return
	}
	//the memory nodes value is similarly presented as CPUset so it can be parsed as cpuset
	val := strings.TrimSpace(memStr)
	nodeList, err := cpuset.Parse(val)
	if err != nil {
		output.IsAligned = false
		output.NNode = -1
		output.Err = fmt.Errorf("could not parse memory nodes' list: %v", err)
		return
	}

	output.ProccessResources.Memory = val
	if nodeList.Size() > 1 {
		output.IsAligned = false
		output.NNode = -1
		return
	}

	node, _ := strconv.Atoi(val)
	if node != output.NNode && output.NNode != -1 {
		output.IsAligned = false
		output.NNode = -1
		return
	}
	output.NNode = node
}
