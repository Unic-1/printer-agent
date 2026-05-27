//go:build windows

package printer

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"printer-agent/models"
)

// discoveredUSB is the JSON shape returned by the PowerShell discovery script.
type discoveredUSB struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

// DiscoverUSB lists USB ESC/POS printers on Windows.
//
// Most thermal printers appear in Device Manager as a printer (USB001, etc.),
// not as a virtual COM port. We therefore enumerate:
//   - Win32_SerialPort (COMx) excluding Bluetooth
//   - Win32_Printer with USB/DOT4/EUSB ports
//   - Installed local printers (RAW via spooler) as a fallback
func DiscoverUSB() ([]*models.Printer, error) {
	found, err := discoverUSBViaPowerShell()
	if err != nil {
		// Fallback for older systems where PowerShell CIM still works but script failed.
		serial, wmicErr := discoverUSBSerialWMIC()
		if wmicErr != nil {
			return nil, fmt.Errorf("usb discovery failed (powershell: %v, wmic: %v)", err, wmicErr)
		}
		return serial, nil
	}
	if len(found) == 0 {
		// Last resort: wmic serial ports only.
		serial, _ := discoverUSBSerialWMIC()
		return serial, nil
	}
	return found, nil
}

func discoverUSBViaPowerShell() ([]*models.Printer, error) {
	script := `
$ErrorActionPreference = 'SilentlyContinue'
$seen = @{}
$out = [System.Collections.ArrayList]@()

function Add-Device($id, $name, $address) {
  if ([string]::IsNullOrWhiteSpace($address)) { return }
  $key = $address.ToLower()
  if ($seen.ContainsKey($key)) { return }
  $seen[$key] = $true
  [void]$out.Add([pscustomobject]@{
    id = $id
    name = $name
    address = $address
  })
}

# 1) Virtual COM ports (USB-serial adapters, some printer drivers)
Get-CimInstance Win32_SerialPort | ForEach-Object {
  $name = $_.Name
  if ($name -match 'Bluetooth') { return }
  $port = $_.DeviceID
  if ($port -match '^COM\d+$') {
    Add-Device ('usb:' + $port) $name $port
  }
}

# 2) Printers exposed on USB/DOT4/EUSB ports (common for thermal receipt printers)
Get-CimInstance Win32_Printer | ForEach-Object {
  $port = $_.PortName
  $pname = $_.Name
  if ([string]::IsNullOrWhiteSpace($port)) { return }
  if ($port -match '^(USB|DOT4|EUSB|LPT)') {
    Add-Device ('usb:port:' + $port) ($pname + ' (' + $port + ')') $port
  }
}

# 3) Local installed printers — print RAW through the Windows spooler by name
Get-CimInstance Win32_Printer | ForEach-Object {
  if (-not $_.Local) { return }
  $pname = $_.Name
  if ([string]::IsNullOrWhiteSpace($pname)) { return }
  $addr = 'printer:' + $pname
  Add-Device ('usb:printer:' + $pname) $pname $addr
}

$out | ConvertTo-Json -Compress
`
	cmd := exec.Command(
		"powershell.exe",
		"-NoProfile",
		"-NonInteractive",
		"-ExecutionPolicy", "Bypass",
		"-Command", script,
	)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("powershell: %w (%s)", err, strings.TrimSpace(string(out)))
	}

	raw := strings.TrimSpace(string(out))
	if raw == "" || raw == "null" {
		return nil, nil
	}

	var rows []discoveredUSB
	if err := json.Unmarshal([]byte(raw), &rows); err != nil {
		// ConvertTo-Json returns a single object when count==1.
		var one discoveredUSB
		if err2 := json.Unmarshal([]byte(raw), &one); err2 != nil {
			return nil, fmt.Errorf("parse discovery json: %w", err)
		}
		rows = []discoveredUSB{one}
	}

	printers := make([]*models.Printer, 0, len(rows))
	for _, row := range rows {
		if row.Address == "" {
			continue
		}
		id := row.ID
		if id == "" {
			id = "usb:" + row.Address
		}
		name := row.Name
		if name == "" {
			name = row.Address
		}
		printers = append(printers, &models.Printer{
			ID:      id,
			Name:    name,
			Type:    models.PrinterUSB,
			Address: row.Address,
			Online:  true,
		})
	}
	return printers, nil
}

func discoverUSBSerialWMIC() ([]*models.Printer, error) {
	cmd := exec.Command("wmic", "path", "Win32_SerialPort", "get", "DeviceID,Name")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	var printers []*models.Printer

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "DeviceID") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		comPort := parts[0]
		name := strings.Join(parts[1:], " ")
		if strings.Contains(strings.ToLower(name), "bluetooth") {
			continue
		}

		printers = append(printers, &models.Printer{
			ID:      "usb:" + comPort,
			Name:    name,
			Type:    models.PrinterUSB,
			Address: comPort,
			Online:  true,
		})
	}

	return printers, nil
}
