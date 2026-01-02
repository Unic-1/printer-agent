package printer

import (
	"os"
)

func printUSB(device string, data []byte) error {
	f, err := os.OpenFile(device, os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	return err
}
