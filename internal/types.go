package internal

import "k8s.io/kubernetes/pkg/kubelet/cm/cpuset"

type NumaAlignmentOutput struct {
	IsAligned         bool
	NNode             int
	Err               error
	ProccessResources ProccessResources
}

type ProccessResources struct {
	CPUs   cpuset.CPUSet
	PCI    []string
	Memory []string
}
