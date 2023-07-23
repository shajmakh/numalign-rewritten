//+++

package main

import (
	"fmt"
	"log"
	"os"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	. "github.com/shajmakh/numaalign-rewritten/internal/cpu"

	"github.com/shajmakh/numaalign-rewritten/pkg/numa"
)

func main() {
	fmt.Println("Verify resources alignment to NUMAs")
	//allow tackling optional flags
	//get resources consumed by the workload
	//get numa locality
	//compare resources alignment to NUMAs
	output := new(NumaAlignmentOutput)

	nNodeCount, err := numa.GetNumaCount()
	if err != nil {
		log.Fatal(err) //TODO vs .Fatalf("%v",err)
	}

	log.Printf("Numa count on system is: %d", nNodeCount)
	if nNodeCount == 1 {
		LogNumaAlignment(NumaAlignmentOutput{
			IsAligned: true,
			NNode:     0,
			Err:       nil,
		})
		os.Exit(0)
	}

	//fail as early as possible and check cpu alignment first
	CheckCpuAlignment(output) //TODO make the option to do the check for other processes based on the provided args/flags
	LogNumaAlignment(*output)

	if false {
		os.Exit(-1)
	}

	os.Exit(0)
}
