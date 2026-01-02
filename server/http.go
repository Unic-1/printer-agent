package server

import "fmt"
import (
	"encoding/json"
	"net/http"
	"printer-agent/models"
	"printer-agent/printer"
)

func Start() {
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/printers", func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode(printer.GetPrinters())
	})


	http.HandleFunc("/printers/register", func(w http.ResponseWriter, r *http.Request) {
    var p models.Printer
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        http.Error(w, err.Error(), 400)
        return
    }

    printer.RegisterPrinter(&p)
    w.Write([]byte("registered"))
	})

	http.HandleFunc("/print", func(w http.ResponseWriter, r *http.Request) {
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
	})

	http.HandleFunc("/bluetooth/devices", func(w http.ResponseWriter, _ *http.Request) {
		devices, err := printer.DiscoverBluetooth()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		json.NewEncoder(w).Encode(devices)
	})


	http.ListenAndServe("127.0.0.1:9123", nil)
}
