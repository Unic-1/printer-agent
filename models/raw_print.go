package models

type RawPrintRequest struct {
	PrinterID string `json:"printerId"`
	Data      string `json:"data"` // base64
}
