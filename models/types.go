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
	PrinterUSB       PrinterType = "usb"
	PrinterBluetooth PrinterType = "bluetooth"
)

// Printer represents a registered printer.
// The frontend may send the printer type as either "type" or "interface" —
// both fields are accepted. The address may contain a "bt:" prefix which
// is stripped before use.
type Printer struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Type      PrinterType `json:"type"`
	Interface PrinterType `json:"interface"` // alias for Type, sent by some clients
	Address   string      `json:"address"`
	Online    bool        `json:"online"`
}

// UnmarshalJSON ensures that if "type" is empty but "interface" is set,
// "interface" is used as the printer type. Also strips "bt:" prefix from address.
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

	// Strip "bt:" prefix from BT device paths (e.g. "bt:/dev/cu.Printer001" → "/dev/cu.Printer001")
	if strings.HasPrefix(p.Address, "bt:") {
		p.Address = strings.TrimPrefix(p.Address, "bt:")
	}

	return nil
}

type PrintRequest struct {
	PrinterID string `json:"printerId"`
	Content   string `json:"content"`
	Cut       bool   `json:"cut"`
}

type BluetoothPrinter struct {
	Path string
	File *os.File
	Mu   sync.Mutex
}