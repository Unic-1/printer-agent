//go:build windows

package printer

import (
	"os"
	"strings"
)

func printUSB(device string, data []byte) error {
	addr := NormalizeUSBAddress(device)
	if strings.HasPrefix(addr, "printer:") {
		return printWindowsSpoolerRAW(strings.TrimPrefix(addr, "printer:"), data)
	}
	return printUSBPort(addr, data)
}

func printUSBPort(device string, data []byte) error {
	path := device
	if !strings.HasPrefix(path, `\\.\`) {
		path = `\\.\` + path
	}

	f, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}
