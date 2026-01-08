//go:build darwin

package printer

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"os"
	"time"
	"printer-agent/models"
)

var printersConnection = map[string]*models.BluetoothPrinter{}

func OpenBluetoothPrinter(path string) (*models.BluetoothPrinter, error) {
	log.Println("Open bluettoth printer " + path)
	// Disable hangup-on-close (HUPCL) BEFORE opening
	// _ = exec.Command("stty", "-f", path, "hupcl").Run()
	// time.Sleep(300 * time.Millisecond)

	// // Disable hangup-on-close for runtime stability
	// if err := exec.Command("stty", "-f", path, "-hupcl").Run(); err != nil {
	// 	return nil, fmt.Errorf("failed to disable hupcl: %w", err)
	// }

	f, err := os.OpenFile(path, os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	// Let RFCOMM stabilize
	time.Sleep(800 * time.Millisecond)

	return &models.BluetoothPrinter{
		Path: path,
		File: f,
	}, nil
}

func PrinterBluetooth(p *models.BluetoothPrinter, data []byte) error {
	fmt.Print("Print via bluetooth", p.Path)
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if p.File == nil {
		return fmt.Errorf("printer not open")
	}

	// Ensure ESC/POS job termination
	if !bytes.HasSuffix(data, []byte("\n\n")) {
		fmt.Print("Add suffix")
		data = append(data, '\n', '\n')
	}

	n, err := p.File.Write(data)
	if err != nil {
		return err
	}

	fmt.Print("Write successful")

	if n != len(data) {
		return fmt.Errorf("partial write: %d/%d", n, len(data))
	}

	return nil
}

func Close(p *models.BluetoothPrinter) error {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if p.File == nil {
		return nil
	}

	time.Sleep(1 * time.Second)
	err := p.File.Close()
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

	_, ok := printersConnection[address]
	if !ok {
		printer, err := OpenBluetoothPrinter(address)
		if err != nil {
			log.Fatal(err)
		}
		printersConnection[address] = printer
	}

	printer, _ := printersConnection[address]

	err := PrinterBluetooth(printer, data)

	time.Sleep(700 * time.Millisecond)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
