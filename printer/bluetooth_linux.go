//go:build linux

package printer

import (
	"os/exec"
	"strings"
	"os"
	"printer-agent/models"
)

func DiscoverBluetooth() ([]*models.Printer, error) {
	cmd := exec.Command("bluetoothctl", "paired-devices")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var printers []*models.Printer

	for _, line := range lines {
		// Example:
		// Device DC:0D:30:7C:22:4E Thermal_Printer
		if !strings.HasPrefix(line, "Device") {
			continue
		}

		parts := strings.SplitN(line, " ", 3)
		if len(parts) < 3 {
			continue
		}

		address := parts[1]
		name := parts[2]

		printers = append(printers, &models.Printer{
			ID:      "bt:" + address,
			Name:    name,
			Type:    models.PrinterBluetooth,
			Address: address,
			Online:  true,
		})
	}

	return printers, nil
}

func printBluetooth(address string, data []byte) error {
	f, err := os.OpenFile(address, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}