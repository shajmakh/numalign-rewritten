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
		1: getCpuset("0-2,6-7"),
	}

	testCases := []struct {
		testMap           map[int]cpuset.CPUSet
		cpuset            cpuset.CPUSet
		expectedNuma      int
		expectedIsAligned bool
	}{
		{numaCpuMap, getCpuset("0-2"), 1, true},
		{numaCpuMap, getCpuset("5"), 0, true},
		{numaCpuMap, getCpuset("3-5,11,13"), 0, true},  //negative
		{numaCpuMap, getCpuset("1,5,9,12"), -1, false}, //negative
		{numaCpuMap, getCpuset("0-2,1,5"), -1, false},  //negative
	}

	for _, c := range testCases {
		out := NewOutput()
		CheckNumaCpuMapping(c.testMap, c.cpuset, &out)
		if out.NNode != c.expectedNuma || out.IsAligned != c.expectedIsAligned {
			t.Fatalf("expected alignment: %t:%d ; actual: %t/%d ; CPU set: [%v]", c.expectedIsAligned, c.expectedNuma, out.IsAligned, out.NNode, c.cpuset)
		}
	}
}

func getCpuset(set string) cpuset.CPUSet {
	cpuset, _ := cpuset.Parse(set)
	return cpuset
}
