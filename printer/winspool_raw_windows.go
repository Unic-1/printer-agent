//go:build windows

package printer

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modWinspool      = syscall.NewLazyDLL("winspool.drv")
	procOpenPrinterW = modWinspool.NewProc("OpenPrinterW")
	procClosePrinter = modWinspool.NewProc("ClosePrinter")
	procStartDocPrinterW = modWinspool.NewProc("StartDocPrinterW")
	procEndDocPrinter    = modWinspool.NewProc("EndDocPrinter")
	procStartPagePrinter = modWinspool.NewProc("StartPagePrinter")
	procEndPagePrinter   = modWinspool.NewProc("EndPagePrinter")
	procWritePrinter     = modWinspool.NewProc("WritePrinter")
)

type docInfo1 struct {
	DocName    *uint16
	OutputFile *uint16
	Datatype   *uint16
}

// printWindowsSpoolerRAW sends ESC/POS bytes to an installed Windows printer by name.
// The printer driver must accept RAW jobs (typical for "Generic / Text Only" thermal drivers).
func printWindowsSpoolerRAW(printerName string, data []byte) error {
	if printerName == "" {
		return fmt.Errorf("empty printer name")
	}

	namePtr, err := syscall.UTF16PtrFromString(printerName)
	if err != nil {
		return err
	}

	var h syscall.Handle
	ret, _, callErr := procOpenPrinterW.Call(
		uintptr(unsafe.Pointer(namePtr)),
		uintptr(unsafe.Pointer(&h)),
		0,
	)
	if ret == 0 {
		return fmt.Errorf("OpenPrinter(%q): %v", printerName, callErr)
	}
	defer procClosePrinter.Call(uintptr(h))

	docName, _ := syscall.UTF16PtrFromString("ESC/POS")
	rawType, _ := syscall.UTF16PtrFromString("RAW")
	di := docInfo1{
		DocName:  docName,
		Datatype: rawType,
	}

	docID, _, callErr := procStartDocPrinterW.Call(
		uintptr(h),
		1,
		uintptr(unsafe.Pointer(&di)),
	)
	if docID == 0 {
		return fmt.Errorf("StartDocPrinter(%q): %v", printerName, callErr)
	}
	defer procEndDocPrinter.Call(uintptr(h))

	procStartPagePrinter.Call(uintptr(h))

	var written uint32
	buf := data
	for len(buf) > 0 {
		chunk := buf
		const maxChunk = 4096
		if len(chunk) > maxChunk {
			chunk = chunk[:maxChunk]
		}
		n, _, callErr := procWritePrinter.Call(
			uintptr(h),
			uintptr(unsafe.Pointer(&chunk[0])),
			uintptr(len(chunk)),
			uintptr(unsafe.Pointer(&written)),
		)
		if n == 0 {
			procEndPagePrinter.Call(uintptr(h))
			return fmt.Errorf("WritePrinter(%q): %v", printerName, callErr)
		}
		buf = buf[len(chunk):]
	}

	procEndPagePrinter.Call(uintptr(h))
	return nil
}
