//go:build darwin

package printer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func printUSB(device string, data []byte) error {
	addr := NormalizeUSBAddress(device)

	if strings.HasPrefix(addr, "printer:") {
		return printCUPSRaw(strings.TrimPrefix(addr, "printer:"), data)
	}

	// Direct device file write for serial paths (/dev/cu.usbserial-*, etc.)
	if strings.HasPrefix(addr, "/dev/") {
		return printDeviceFile(addr, data)
	}

	// If it looks like a CUPS printer name (no slashes, no dev path), try CUPS
	if !strings.Contains(addr, "/") && addr != "" {
		return printCUPSRaw(addr, data)
	}

	return printDeviceFile(addr, data)
}

func printDeviceFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("open device %s: %w", path, err)
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}

// printCUPSRaw sends raw ESC/POS bytes to a CUPS-registered printer.
// Equivalent to: lp -d <printer> -o raw <tempfile>
func printCUPSRaw(printerName string, data []byte) error {
	if printerName == "" {
		return fmt.Errorf("empty printer name")
	}

	tmpFile, err := os.CreateTemp("", "escpos-*.bin")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("write temp file: %w", err)
	}
	tmpFile.Close()

	var stderr bytes.Buffer
	cmd := exec.Command("lp", "-d", printerName, "-o", "raw", tmpFile.Name())
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg != "" {
			return fmt.Errorf("lp -d %s: %s", printerName, errMsg)
		}
		return fmt.Errorf("lp -d %s: %w", printerName, err)
	}

	return nil
}
