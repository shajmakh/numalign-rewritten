package device

import (
	"fmt"
	"os"
	"strings"

	. "github.com/shajmakh/numaalign-rewritten/internal"
	. "github.com/shajmakh/numaalign-rewritten/pkg/numa"
)

// CheckPciDevicesAlignment checks alignment to numa node of the PCI devices used by the process. The expected used devices are fetched from environment variable DEV_RESOURCES
func CheckPciDevicesAlignment(out *NumaAlignmentOutput) {
	requestedDevs := parseDevicesFromEnv()
	if len(requestedDevs) == 0 {
		return
	}

	deviceNumaMap, err := GetNumaDeviceMapping()
	if err != nil {
		out.Err = err
		return
	}
	for _, devName := range requestedDevs {
		if nnode, found := deviceNumaMap[devName]; found {
			if nnode != out.NNode {
				out.NNode = -1
				if Verbose {
					WriteToDest(fmt.Sprintf("resources used by the process are not aligned to a single numa: PCI device %q\n", devName))
				}
				return
			}
		}
	}
}

func parseDevicesFromEnv() []string {
	devStr, ok := os.LookupEnv("DEV_RESOURCES")
	if !ok {
		if Verbose {
			WriteToDest("no pci devices used")
		}
		return []string{}
	}
	listStr := strings.ReplaceAll(devStr, " ", "")
	return strings.Split(listStr, ",")
}
