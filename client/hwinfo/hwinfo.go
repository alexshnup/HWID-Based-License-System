package hwinfo

import (
	"fmt"

	"github.com/jaypipes/ghw"
)

// HardwareInfo contains relevant hardware information
type HardwareInfo struct {
	Block *ghw.BlockInfo
	Disk  string
}

// GetHardwareInfo initializes and retrieves hardware information
func GetHardwareInfo() (*HardwareInfo, error) {
	block, err := ghw.Block()
	if err != nil {
		return nil, fmt.Errorf("error getting block storage info: %v", err)
	}

	diskSerial := ""
	if len(block.Disks) > 0 {
		diskSerial = block.Disks[0].SerialNumber
	}

	return &HardwareInfo{
		Block: block,
		Disk:  diskSerial,
	}, nil
}
