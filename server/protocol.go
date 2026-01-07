package server

import (
	"encoding/json"
	"log"
	"net/url"

	"printer-agent/commands"
	"printer-agent/models"
)

func handlePrint(q url.Values) {
	req := models.PrintRequest{
		PrinterID: q.Get("printerId"),
		Cut:       q.Get("cut") == "true",
	}

	if err := json.Unmarshal([]byte(q.Get("content")), &req.Content); err != nil {
		log.Println("Invalid content:", err)
		return
	}

	if err := commands.Print(req); err != nil {
		log.Println("Print failed:", err)
	}
}

func handleRawPrint(q url.Values) {
	req := models.RawPrintRequest{
		PrinterID: q.Get("printerId"),
		Data:      q.Get("data"),
	}

	if err := commands.RawPrint(req); err != nil {
		log.Println("Raw print failed:", err)
	}
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}


func HandleProtocol(raw string) {
	u, err := url.Parse(raw)
	if err != nil {
		log.Println("Invalid protocol URL:", err)
		return
	}

	action := u.Host
	q := u.Query()

	switch action {

	case "health":
		log.Println(commands.Health())

	case "print":
		handlePrint(q)

	case "print-raw":
		handleRawPrint(q)

	case "printers":
		data := commands.ListPrinters()
		log.Println(toJSON(data))

	case "bluetooth":
		devices, err := commands.DiscoverBluetooth()
		if err != nil {
			log.Println("Bluetooth error:", err)
			return
		}
		log.Println(toJSON(devices))

	default:
		log.Println("Unknown action:", action)
	}
}
