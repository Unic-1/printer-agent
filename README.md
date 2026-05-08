# Aharsuchi Printer Agent

Local HTTPS printer agent for Aharsuchi POS. Runs as a **Windows Service** that starts automatically at boot and auto-restarts on failure.

## Quick Start (Development)

```bash
# Run interactively (foreground, for debugging)
go run . run
--- or ---
CERT_PATH=./certs go run main.go
```

## Build for Windows

```bash
# Build the Windows executable (GUI subsystem — no console window)
GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o printer-agent.exe
```

## Windows Service Management

Requires **Administrator** privileges:

```cmd
:: Install the service (registers with Windows SCM)
printer-agent.exe install

:: Start the service
printer-agent.exe start

:: Stop the service
printer-agent.exe stop

:: Uninstall the service
printer-agent.exe uninstall

:: Run interactively for debugging (does NOT use SCM)
printer-agent.exe run
```

After installation, the service:

- Appears in `services.msc` as **"Aharsuchi Printer Agent"**
- Starts automatically on boot (delayed auto-start)
- Auto-restarts within 10 seconds if it crashes
- Listens on `https://127.0.0.1:9123`

## Build Installer (Inno Setup)

```bash
# Build the exe first
GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o printer-agent.exe

# Build installer using Docker
docker run --rm -i -v "$PWD:/work" amake/innosetup installer.iss
```

The installer will:

1. Install the executable and certs to `C:\Program Files\Aharsuchi\`
2. Install the root CA certificate to the Windows trust store
3. Register and start the Windows Service
4. Configure automatic restart on failure

## TLS Certificates

```bash
# Generate certs for localhost
mkcert 127.0.0.1 localhost

# Find the root CA (needed for installer)
mkcert -CAROOT
```

## Cross-platform Builds

```bash
GOOS=windows GOARCH=amd64 go build -ldflags="-H=windowsgui" -o build/printer-agent.exe
GOOS=darwin GOARCH=amd64 go build -o build/printer-agent
GOOS=linux GOARCH=amd64 go build -o build/printer-agent
```

## Troubleshooting

**Mac Bluetooth issue** — if printing fails after process restart:

```bash
sudo pkill bluetoothd
```

**Windows Service won't start** — check Windows Event Viewer:

1. Open Event Viewer (`eventvwr.msc`)
2. Navigate to: Windows Logs → Application
3. Filter by source: `AharsuchiPrinterAgent`
