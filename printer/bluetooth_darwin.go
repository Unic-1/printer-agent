//go:build darwin

package printer

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"printer-agent/models"
)

func OpenBluetoothPrinter(path string) (*models.BluetoothPrinter, error) {
	log.Printf("opening bluetooth printer %s", path)
	resetBluetoothPort(path)

	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	// Let RFCOMM stabilize after open.
	time.Sleep(500 * time.Millisecond)

	return &models.BluetoothPrinter{
		Path: path,
		File: f,
	}, nil
}

func PrinterBluetooth(p *models.BluetoothPrinter, data []byte) error {
	log.Printf("printing via bluetooth %s", p.Path)
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if p.File == nil {
		return fmt.Errorf("printer not open")
	}

	// Ensure ESC/POS job termination
	if !bytes.HasSuffix(data, []byte("\n\n")) {
		data = append(data, '\n', '\n')
	}

	n, err := p.File.Write(data)
	if err != nil {
		return err
	}

	if n != len(data) {
		return fmt.Errorf("partial write: %d/%d", n, len(data))
	}

	if err := drainBluetoothOutput(int(p.File.Fd()), btDrainTimeout); err != nil {
		return fmt.Errorf("output drain failed: %w", err)
	}

	log.Printf("bluetooth write successful on %s", p.Path)
	return nil
}

func Close(p *models.BluetoothPrinter) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if p.File == nil {
		return nil
	}

	err := closeBluetoothPort(p.Path, p.File)
	p.File = nil
	return err
}

func DiscoverBluetooth() ([]*models.Printer, error) {
	// 1. Get paired Bluetooth devices
	cmd := exec.Command("system_profiler", "SPBluetoothDataType")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 2. Read available RFCOMM devices
	devEntries, err := os.ReadDir("/dev")
	if err != nil {
		return nil, err
	}

	var cuDevices []string
	for _, e := range devEntries {
		name := e.Name()
		if strings.HasPrefix(name, "cu.") &&
			!strings.Contains(name, "Bluetooth-Incoming-Port") {
			cuDevices = append(cuDevices, "/dev/"+name)
		}
	}

	lines := strings.Split(string(out), "\n")

	var printers []*models.Printer
	var name string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Device name block
		if strings.HasSuffix(line, ":") &&
			!strings.Contains(line, "Devices") &&
			!strings.Contains(line, "Bluetooth") {
			name = strings.TrimSuffix(line, ":")
			continue
		}

		// We only expose printers that have a cu.* mapping
		if name != "" {
			for _, cu := range cuDevices {
				// loose match: printer name appears in device path
				if strings.Contains(strings.ToLower(cu), strings.ToLower(name)) {
					printers = append(printers, &models.Printer{
						ID:      "bt:" + cu,
						Name:    name,
						Type:    models.PrinterBluetooth,
						Address: cu,
						Online:  true,
					})

					name = ""
					break
				}
			}
		}
	}

	return printers, nil
}

func printBluetooth(address string, data []byte) error {
	var lastErr error
	for cycle := 1; cycle <= 2; cycle++ {
		err := printBluetoothOnce(address, data)
		if err == nil {
			return nil
		}
		lastErr = err
		log.Printf("bluetooth print cycle %d failed for %s: %v", cycle, address, err)
		if cycle < 2 {
			time.Sleep(time.Duration(cycle) * time.Second)
		}
	}
	return fmt.Errorf("print failed on %s: %w", address, lastErr)
}

func printBluetoothOnce(address string, data []byte) error {
	var p *models.BluetoothPrinter
	var openErr error
	for attempt := 1; attempt <= 3; attempt++ {
		p, openErr = OpenBluetoothPrinter(address)
		if openErr == nil {
			break
		}
		log.Printf("open bluetooth printer %s attempt %d failed: %v", address, attempt, openErr)
		if attempt < 3 {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	if openErr != nil {
		return fmt.Errorf("failed to open bluetooth printer %s after retries: %w", address, openErr)
	}
	defer Close(p)

	if err := PrinterBluetooth(p, data); err != nil {
		return err
	}

	time.Sleep(350 * time.Millisecond)
	return nil
}
