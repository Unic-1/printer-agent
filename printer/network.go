package printer

import (
	"net"
	"time"
)

func printNetwork(address string, data []byte) error {
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.SetWriteDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}
