//go:build darwin

package printer

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"golang.org/x/sys/unix"
)

const (
	btPortResetDelay   = 200 * time.Millisecond
	btDrainTimeout     = 2 * time.Second
	btTeardownDelay    = 200 * time.Millisecond
)

// resetBluetoothPort drops any ghost RFCOMM session left by a prior process
// before opening the device for a new print job.
func resetBluetoothPort(path string) {
	cmd := exec.Command("stty", "-f", path, "hupcl")
	if err := cmd.Run(); err != nil {
		log.Printf("warning: reset bluetooth port %s: %v", path, err)
	}
	time.Sleep(btPortResetDelay)
}

// drainBluetoothOutput waits until buffered output is transmitted. A timeout
// indicates a dead or half-open RFCOMM channel.
func drainBluetoothOutput(fd int, timeout time.Duration) error {
	done := make(chan error, 1)
	go func() {
		_, _, errno := unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.TIOCDRAIN), 0)
		if errno != 0 {
			done <- errno
			return
		}
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("tcdrain timed out after %s", timeout)
	}
}

// closeBluetoothPort tears down an RFCOMM session so bluetoothd releases the channel.
func closeBluetoothPort(path string, f *os.File) error {
	if f == nil {
		return nil
	}

	fd := int(f.Fd())
	if err := drainBluetoothOutput(fd, btDrainTimeout); err != nil {
		log.Printf("warning: drain on close for %s: %v", path, err)
	}

	err := f.Close()
	if sttyErr := exec.Command("stty", "-f", path, "hupcl").Run(); sttyErr != nil {
		log.Printf("warning: stty hupcl on %s: %v", path, sttyErr)
	}
	time.Sleep(btTeardownDelay)
	return err
}
