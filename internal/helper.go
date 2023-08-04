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

package internal

import (
	"fmt"
	"io"
	"log"
)

var Verbose bool

// LogNumaAlignment prints the final result of the program, -1 if resources are not aligned,otherwise the numa on which all the resources are aligned
func LogNumaAlignment(res NumaAlignmentOutput, dest io.Writer) {
	WriteToDest(fmt.Sprintf("NUMA %d\n", res.NNode), dest)

	if Verbose {
		printResources(res.ProccessResources)
		if res.Err != nil {
			fmt.Printf("Error: %v\n", res.Err)
		}
	}
}

func printResources(rsrc ProccessResources) { //could be done a ToString() instead but would it be worth it to have another file for the process details (=app output)?
	log.Printf("consumed resources:\n CPUs:\n%v\n PCI devices:\n%v\n Memory:\n%v\n", rsrc.CPUs.String(), rsrc.PCI, rsrc.Memory)
}

func WriteToDest(str string, dest io.Writer) {
	_, err := io.WriteString(dest, str)
	if err != nil {
		log.Fatal(err)
	}
}
