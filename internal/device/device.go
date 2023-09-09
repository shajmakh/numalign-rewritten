package device

import (
	"fmt"
	"os"
	"strings"

	"github.com/shajmakh/numaalign-rewritten/internal"
	"github.com/shajmakh/numaalign-rewritten/pkg/numa"
)

// CheckPciDevicesAlignment checks alignment to numa node of the PCI devices used by the process. The expected used devices are fetched from environment variable DEV_RESOURCES
func CheckPciDevicesAlignment(out *internal.NumaAlignmentOutput) {
	requestedDevs := parseDevicesFromEnv()
	if len(requestedDevs) == 0 {
		return
	}

	deviceNumaMap, err := numa.GetNumaDeviceMapping()
	if err != nil {
		out.Err = err
		return
	}

	CheckPciDeviceToNumaMapping(deviceNumaMap, requestedDevs, out)
}

/*
CheckPciDeviceToNumaMapping updates "out" with the alignment result of device resources to a single numa.
If out.NNode is not -1 it compares the numa of the devices with that numa and if numas are not the same it
updates "out" with un-alignment info --> IsAligned: false; NNode:-1
*/
func CheckPciDeviceToNumaMapping(deviceNumaMap map[string]int, devList []string, out *internal.NumaAlignmentOutput) {
	if !out.IsAligned {
		return
	}
	if len(devList) == 0 {
		return
	}
	for _, devName := range devList {
		if nnode, found := deviceNumaMap[devName]; found {
			if out.NNode == -1 {
				out.NNode = nnode
				continue
			}

			if nnode != out.NNode {
				out.NNode = -1
				out.IsAligned = false
				if internal.Verbose {
					internal.WriteToDest(fmt.Sprintf("resources used by the process are not aligned to a single numa: PCI device %q\n", devName))
				}
				return
			}
		}
	}
}

func parseDevicesFromEnv() []string {
	devStr, ok := os.LookupEnv("DEV_RESOURCES")
	if !ok {
		if internal.Verbose {
			internal.WriteToDest("no pci devices used")
		}
		return []string{}
	}
	listStr := strings.ReplaceAll(devStr, " ", "")
	return strings.Split(listStr, ",")
}
