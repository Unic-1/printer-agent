package models

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

type PrinterType string

const (
	PrinterNetwork   PrinterType = "network"
	PrinterIP        PrinterType = "ip" // alias for PrinterNetwork; frontend may send "ip"
	PrinterUSB       PrinterType = "usb"
	PrinterBluetooth PrinterType = "bluetooth"
)

// Printer represents a registered printer.
// The frontend may send the printer type as either "type" or "interface" —
// both fields are accepted. Bluetooth addresses may use a "bt:" prefix
// (e.g. "bt:/dev/cu.Printer001"); that prefix is kept for registration and
// print lookup and stripped only when opening the device (see BluetoothDevicePath).
type Printer struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      PrinterType `json:"type"`
	Interface PrinterType `json:"interface"` // alias for Type, sent by some clients
	Address   string      `json:"address"`
	Online    bool        `json:"online"`
}

// BluetoothDevicePath returns the OS device path for a Bluetooth printer address,
// stripping the optional "bt:" prefix used by the frontend and discovery API.
func BluetoothDevicePath(address string) string {
	return strings.TrimPrefix(address, "bt:")
}

// UnmarshalJSON ensures that if "type" is empty but "interface" is set,
// "interface" is used as the printer type.
func (p *Printer) UnmarshalJSON(b []byte) error {
	// Use an alias to avoid infinite recursion
	type PrinterAlias Printer
	aux := &PrinterAlias{}
	if err := json.Unmarshal(b, aux); err != nil {
		return err
	}
	*p = Printer(*aux)

	// Resolve type: prefer "type", fall back to "interface"
	if p.Type == "" && p.Interface != "" {
		p.Type = p.Interface
	}

	return nil
}

type PrintRequest struct {
	PrinterID string `json:"printerId"`
	Content   string `json:"content"`
	Cut       bool   `json:"cut"`
}

// RegisteredPrintersResponse is returned by GET /printers and GET /printers/registered.
type RegisteredPrintersResponse struct {
	Count    int        `json:"count"`
	Printers []*Printer `json:"printers"`
}

type BluetoothPrinter struct {
	Path string
	File *os.File
	Mu   sync.Mutex
}