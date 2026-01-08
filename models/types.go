package models

import (
	"os"
	"sync"
)

type PrinterType string

const (
	PrinterNetwork   PrinterType = "network"
	PrinterUSB       PrinterType = "usb"
	PrinterBluetooth PrinterType = "bluetooth"
)

type Printer struct {
	ID      string      `json:"id"`
	Name    string      `json:"name"`
	Type    PrinterType `json:"type"`
	Address string      `json:"address"`
	Online  bool        `json:"online"`
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