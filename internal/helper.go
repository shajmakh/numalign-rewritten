package internal

import (
	"fmt"
	"log"
)

func LogNumaAlignment(res NumaAlignmentOutput) {
	//TODO print better considering the e2e tests
	printResources(res.ProccessResources)
	if !res.IsAligned {
		log.Println("Resources are not aligned to a single numa")
		if res.Err != nil {
			log.Printf("Error: %v", res.Err)
		}
		return
	}
	log.Printf("Resources are aligned to numa %d", res.NNode)
}

func printResources(rsrc ProccessResources) { //could be done a ToString() instead but would it be worth it to have another file for the process details (=app output)?
	fmt.Printf("consumed resources:\n CPUs:\n%v\n PCI devices:\n%v\n Memory:\n%v\n", rsrc.CPUs.String(), rsrc.PCI, rsrc.Memory)
}
