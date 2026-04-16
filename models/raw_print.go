package models

type RawPrintRequest struct {
	PrinterID string `json:"PrinterId"` // frontend sends capital-I
	Data      string `json:"Data"`      // base64
}
