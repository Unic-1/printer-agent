package main

import (
	"printer-agent/models"
	"printer-agent/printer"
	"printer-agent/server"
)

func main() {
	// Example pre-registered printers
	printer.RegisterPrinter(&models.Printer{
		ID:      "net-1",
		Name:    "Kitchen Printer",
		Type:    models.PrinterNetwork,
		Address: "192.168.1.50:9100",
		Online:  true,
	})

	printer.RegisterPrinter(&models.Printer{
		ID:      "usb-1",
		Name:    "Counter Printer",
		Type:    models.PrinterUSB,
		Address: "/dev/usb/lp0",
		Online:  true,
	})

	server.Start()
}
