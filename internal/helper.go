package internal

import (
	"fmt"
	"log"
)

// LogNumaAlignment prints the final result of the program, -1 if resources are not aligned,otherwise the numa on which all the resources are aligned
func LogNumaAlignment(res NumaAlignmentOutput) {
	fmt.Printf("NUMA %d\n", res.NNode)
	//TODO make the below a debug level info
	printResources(res.ProccessResources)
	if res.Err != nil {
		fmt.Printf("Error: %v", res.Err)
	}
}

func printResources(rsrc ProccessResources) { //could be done a ToString() instead but would it be worth it to have another file for the process details (=app output)?
	log.Printf("consumed resources:\n CPUs:\n%v\n PCI devices:\n%v\n Memory:\n%v\n", rsrc.CPUs.String(), rsrc.PCI, rsrc.Memory)
}
