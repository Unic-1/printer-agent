package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kardianos/service"

	"printer-agent/printer"
	"printer-agent/server"
)

// program implements the kardianos/service Interface.
// It manages the HTTP/TLS server lifecycle.
type program struct {
	srv      *http.Server
	certFile string
	keyFile  string
	logger   service.Logger
}

func (p *program) Start(s service.Service) error {
	p.logger.Info("Starting Aharsuchi Printer Agent...")

	// Resolve cert directory: CERT_PATH env overrides the default next to the executable.
	// Default is required because Windows Services run with CWD = C:\Windows\System32.
	certDir := os.Getenv("CERT_PATH")
	if certDir == "" {
		certDir = filepath.Join(server.ExeDir(), "certs")
	}

	p.srv, p.certFile, p.keyFile = server.NewServer(certDir)

	// Start the HTTPS server in a background goroutine.
	// service.Start() must return quickly — the SCM will kill us if we block.
	go func() {
		p.logger.Infof("Printer agent running on https://127.0.0.1:9123")
		p.logger.Infof("Using certs from: %s", certDir)

		err := p.srv.ListenAndServeTLS(p.certFile, p.keyFile)
		if err != nil && err != http.ErrServerClosed {
			p.logger.Errorf("HTTPS server error: %v", err)
			// Exit non-zero so Windows SCM triggers a restart via recovery actions
			os.Exit(1)
		}
	}()

	return nil
}

func (p *program) Stop(s service.Service) error {
	p.logger.Info("Stopping Aharsuchi Printer Agent...")

	printer.Shutdown()

	if p.srv != nil {
		// Give in-flight requests up to 5 seconds to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return p.srv.Shutdown(ctx)
	}
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "AharsuchiPrinterAgent",
		DisplayName: "Aharsuchi Printer Agent",
		Description: "Local HTTPS print server for Aharsuchi POS – manages network, USB, and Bluetooth receipt printers.",
		Option: service.KeyValue{
			// Delayed auto-start: let Windows finish booting before starting us.
			// This avoids race conditions with Bluetooth/network drivers.
			"DelayedAutoStart": true,
		},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Obtain a logger that writes to the Windows Event Log (or syslog on Linux/mac).
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	prg.logger = logger

	// --- CLI command handling ---
	// Usage:
	//   printer-agent.exe install    — register as a Windows Service
	//   printer-agent.exe uninstall  — remove the Windows Service
	//   printer-agent.exe start      — start the service via SCM
	//   printer-agent.exe stop       — stop the service via SCM
	//   printer-agent.exe run        — run interactively (foreground, for debugging)
	//   printer-agent.exe            — run as service (called by SCM)

	if len(os.Args) > 1 {
		cmd := os.Args[1]
		switch cmd {
		case "install":
			err = s.Install()
			if err != nil {
				log.Fatalf("Failed to install service: %v", err)
			}
			fmt.Println("Service installed successfully.")
			fmt.Println("Run 'printer-agent.exe start' to start the service.")
			return

		case "uninstall":
			// Stop the service first (ignore errors — it may not be running)
			_ = s.Stop()
			err = s.Uninstall()
			if err != nil {
				log.Fatalf("Failed to uninstall service: %v", err)
			}
			fmt.Println("Service uninstalled successfully.")
			return

		case "start":
			err = s.Start()
			if err != nil {
				log.Fatalf("Failed to start service: %v", err)
			}
			fmt.Println("Service started.")
			return

		case "stop":
			err = s.Stop()
			if err != nil {
				log.Fatalf("Failed to stop service: %v", err)
			}
			fmt.Println("Service stopped.")
			return

		case "run":
			// Run interactively in the foreground (for debugging).
			// This does NOT register with the SCM — just runs the program.
			fmt.Println("Running in interactive mode (Ctrl+C to stop)...")
			err = s.Run()
			if err != nil {
				logger.Error(err)
			}
			return

		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
			fmt.Fprintln(os.Stderr, "Usage: printer-agent.exe [install|uninstall|start|stop|run]")
			os.Exit(1)
		}
	}

	// No arguments → run as service (invoked by Windows SCM).
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
