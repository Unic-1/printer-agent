package printer

import (
	"errors"
	"strings"

	"printer-agent/models"
)

var printers = map[string]*models.Printer{}

func printerKey(p *models.Printer) string {
	if p.ID != "" {
		return p.ID
	}
	return p.Address
}

// lookupKeys returns printer map keys to try for a print request id.
// Clients register and print with "bt:/dev/cu.*" while older entries may
// only be stored under the bare device path.
func lookupKeys(printerID string) []string {
	keys := []string{printerID}
	if strings.HasPrefix(printerID, "bt:") {
		keys = append(keys, strings.TrimPrefix(printerID, "bt:"))
	} else if strings.HasPrefix(printerID, "usb:") {
		keys = append(keys, strings.TrimPrefix(printerID, "usb:"))
	} else if printerID != "" {
		keys = append(keys, "bt:"+printerID, "usb:"+printerID)
	}
	return keys
}

func resolvePrinter(printerID string) (*models.Printer, bool) {
	for _, key := range lookupKeys(printerID) {
		if p, ok := printers[key]; ok {
			return p, true
		}
	}
	return nil, false
}

func RegisterPrinter(p *models.Printer) {
	if p.ID == "" && p.Address != "" {
		p.ID = p.Address
	}
	printers[printerKey(p)] = p
}

// ListRegistered returns all printers currently registered with the agent.
func ListRegistered() models.RegisteredPrintersResponse {
	list := make([]*models.Printer, 0, len(printers))
	for _, p := range printers {
		list = append(list, p)
	}
	return models.RegisteredPrintersResponse{
		Count:    len(list),
		Printers: list,
	}
}

func Print(printerID string, data []byte) error {
	p, ok := resolvePrinter(printerID)
	if !ok {
		return errors.New("printer not found")
	}

	switch p.Type {
	case models.PrinterNetwork, models.PrinterIP:
		// Ensure a port is present; ESC/POS printers default to 9100.
		return printNetwork(ensureNetworkPort(p.Address), data)
	case models.PrinterUSB:
		addr := NormalizeUSBAddress(p.Address)
		if addr == "" {
			addr = NormalizeUSBAddress(p.ID)
		}
		return printUSB(addr, data)
	case models.PrinterBluetooth:
		return printBluetooth(models.BluetoothDevicePath(p.Address), data)
	default:
		return errors.New("unsupported printer type")
	}
}

// ensureNetworkPort appends the default ESC/POS port (9100) when the address
// does not already contain a port number.
func ensureNetworkPort(address string) string {
	if strings.Contains(address, ":") {
		return address
	}
	return address + ":9100"
}
