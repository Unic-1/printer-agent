//go:build !windows && !darwin && !linux

package printer

import "errors"

func printBluetooth(address string, data []byte) error {
	return errors.New("bluetooth printing not enabled in this build")
}
