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
	"log"
	"os"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	. "github.com/shajmakh/numaalign-rewritten/internal/cpu"

	"github.com/shajmakh/numaalign-rewritten/pkg/numa"
)

func main() {
	//TODO allow optional flags like pid, output file,debug level

	output := new(NumaAlignmentOutput)

	nNodeCount, err := numa.GetNumaCount()
	if err != nil {
		log.Fatal(err) //TODO vs .Fatalf("%v",err)
	}
	log.Printf("Numa count on system is: %d", nNodeCount) //TODO make it a debug level output
	if nNodeCount == 1 {
		LogNumaAlignment(NumaAlignmentOutput{
			NNode: 0,
			Err:   nil,
		})
		os.Exit(0)
	}

	CheckCpuAlignment(output)

	LogNumaAlignment(*output)

	if false {
		os.Exit(-1)
	}

	os.Exit(0)
}
