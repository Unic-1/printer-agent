package server

import "fmt"

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"printer-agent/models"
	"printer-agent/printer"
)

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func listPrinters(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(printer.GetPrinters())
}

func registerPrinter(w http.ResponseWriter, r *http.Request) {
	var p models.Printer
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), 400)
			return
	}

	printer.RegisterPrinter(&p)
	w.Write([]byte("registered"))
}

func print(w http.ResponseWriter, r *http.Request) {
	var req models.PrintRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Println("Printer ID:")
	fmt.Println(req.PrinterID)

	data := printer.BuildEscPos(req.Content, req.Cut)

	if err := printer.Print(req.PrinterID, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Write([]byte("printed"))
}

func rawPrint(w http.ResponseWriter, r *http.Request) {
	var req models.RawPrintRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.PrinterID == "" || req.Data == "" {
		http.Error(w, "printerId and data are required", http.StatusBadRequest)
		return
	}

	// Decode base64 → []byte
	raw, err := base64.StdEncoding.DecodeString(req.Data)
	if err != nil {
		http.Error(w, "invalid base64 data", http.StatusBadRequest)
		return
	}

	// Directly write bytes to printer
	if err := printer.Print(req.PrinterID, raw); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("printed"))
}

func deviceList(w http.ResponseWriter, r *http.Request) {
	devices, err := printer.DiscoverBluetooth()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(devices)
}


func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", health)
	mux.HandleFunc("/printers", listPrinters)
	mux.HandleFunc("/printers/register", registerPrinter)
	mux.HandleFunc("/print", print)
	mux.HandleFunc("/print/raw", rawPrint)
	mux.HandleFunc("/bluetooth/devices", deviceList)

	http.ListenAndServe("127.0.0.1:9123", withCORS(mux))
}
