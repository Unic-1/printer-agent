//go:build windows

package printer

import (
	"os/exec"
	"strings"
	"os"
	"printer-agent/models"
)

func DiscoverBluetooth() ([]*models.Printer, error) {
	cmd := exec.Command("wmic", "path", "Win32_SerialPort", "get", "DeviceID,Name")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var printers []*models.Printer

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "DeviceID") {
			continue
		}

		// Example:
		// COM5  Standard Serial over Bluetooth link
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		comPort := parts[0]
		name := strings.Join(parts[1:], " ")

		// Filter Bluetooth ports only
		if !strings.Contains(strings.ToLower(name), "bluetooth") {
			continue
		}

		printers = append(printers, &models.Printer{
			ID:      "bt:" + comPort,
			Name:    name,
			Type:    models.PrinterBluetooth,
			Address: comPort,
			Online:  true,
		})
	}

	return printers, nil
}

func printBluetooth(address string, data []byte) error {
	port := `\\.\` + address // Required Windows syntax

	f, err := os.OpenFile(port, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}