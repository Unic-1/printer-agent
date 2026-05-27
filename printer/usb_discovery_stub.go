//go:build !darwin && !linux && !windows

package printer

import "printer-agent/models"

func DiscoverUSB() ([]*models.Printer, error) {
	return nil, nil
}
