//go:build darwin

package printer

import (
	"os/exec"
	"strings"
	"os"
	"time"
	"printer-agent/models"
)

func DiscoverBluetooth() ([]*models.Printer, error) {
	cmd := exec.Command("system_profiler", "SPBluetoothDataType")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var printers []*models.Printer

	var name, address string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasSuffix(line, ":") && !strings.Contains(line, "Devices") {
			name = strings.TrimSuffix(line, ":")
		}

		if strings.HasPrefix(line, "Address:") {
			address = strings.TrimSpace(strings.Replace(line, "Address:", "", 1))

			printers = append(printers, &models.Printer{
				ID:      "bt:" + address,
				Name:    name,
				Type:    models.PrinterBluetooth,
				Address: address,
				Online:  true,
			})

			name = ""
			address = ""
		}
	}

	return printers, nil
}

func printBluetooth(address string, data []byte) error {
	// IMPORTANT: must be O_RDWR on macOS
	f, err := os.OpenFile(address, os.O_RDWR, 0666)
	if err != nil {
		return err
	}

	// Give macOS time to establish RFCOMM connection
	time.Sleep(300 * time.Millisecond)

	_, err = f.Write(data)

	// Ensure data is flushed
	f.Sync()
	f.Close()

	// Give printer time to process
	time.Sleep(200 * time.Millisecond)

	return err
}