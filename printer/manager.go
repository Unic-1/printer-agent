package printer

import (
	"errors"
	"printer-agent/models"
)

var printers = map[string]*models.Printer{}

func RegisterPrinter(p *models.Printer) {
	printers[p.ID] = p
}

func GetPrinters() []*models.Printer {
	list := []*models.Printer{}
	for _, p := range printers {
		list = append(list, p)
	}
	return list
}

func Print(printerID string, data []byte) error {
	p, ok := printers[printerID]
	if !ok {
		return errors.New("printer not found")
	}

	switch p.Type {
	case models.PrinterNetwork:
		return printNetwork(p.Address, data)
	case models.PrinterUSB:
		return printUSB(p.Address, data)
	case models.PrinterBluetooth:
		return printBluetooth(p.Address, data)
	default:
		return errors.New("unsupported printer type")
	}
}
