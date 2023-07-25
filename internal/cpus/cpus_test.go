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
		{numaCpuMap, getCpuset("1,5,9,12"), -1}, //negative
		{numaCpuMap, getCpuset("0-2,1,5"), -1},  //negative
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
