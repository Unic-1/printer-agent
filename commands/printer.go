package commands

import (
	"encoding/base64"
	"errors"

	"printer-agent/models"
	"printer-agent/printer"
)

func Health() string {
	return "ok"
}

func ListPrinters() interface{} {
	return printer.GetPrinters()
}

func RegisterPrinter(p *models.Printer) error {
	printer.RegisterPrinter(p)
	return nil
}

func Print(req models.PrintRequest) error {
	data := printer.BuildEscPos(req.Content, req.Cut)
	return printer.Print(req.PrinterID, data)
}

func RawPrint(req models.RawPrintRequest) error {
	if req.PrinterID == "" || req.Data == "" {
		return errors.New("printerId and data are required")
	}

	raw, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		return errors.New("invalid base64 data")
	}

	return printer.Print(req.PrinterID, raw)
}

func DiscoverBluetooth() (interface{}, error) {
	return printer.DiscoverBluetooth()
}
