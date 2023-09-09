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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/shajmakh/numaalign-rewritten/internal"
	"github.com/shajmakh/numaalign-rewritten/internal/cpu"
	"github.com/shajmakh/numaalign-rewritten/internal/device"
	"github.com/shajmakh/numaalign-rewritten/internal/memory"

	"github.com/shajmakh/numaalign-rewritten/pkg/numa"
)

var (
	pid            = flag.String("p", "", "pid of the process for which to check the numa alignment of its resources")
	verbose        = flag.Bool("v", false, "display app output with debug level")
	outputFilePath = flag.String("o", "", "path of output file; leave empty to display on standart output")
)

func main() {
	flag.Parse()

	if outputFilePath != nil && *outputFilePath != "" {
		f, err := os.Create(*outputFilePath)
		if err != nil {
			log.Fatalf("error opening %s: %v\n", *outputFilePath, err)
		}
		defer f.Close()
		internal.OutputDest = f
	}

	internal.Verbose = *verbose

	processId := "self"
	if strings.TrimSpace(*pid) != "" {
		processId = strings.Fields(*pid)[0]
	}

	output := internal.NewOutput()

	nNodeCount, err := numa.GetNumaCount()
	if err != nil {
		log.Fatal(err) //TODO vs .Fatalf("%v",err)
	}

	if internal.Verbose {
		internal.WriteToDest(fmt.Sprintf("Numa count on system is: %d", nNodeCount))
	}

	if nNodeCount == 1 {
		output.NNode = 0
		internal.LogNumaAlignment(output)
		os.Exit(0)
	}

	cpu.CheckCpuAlignment(processId, &output)

	device.CheckPciDevicesAlignment(&output)

	memory.CheckMemoryResourcesAlignment(processId, &output)

	internal.LogNumaAlignment(output)

	if !output.IsAligned {
		os.Exit(-1)
	}

	os.Exit(0)
}
