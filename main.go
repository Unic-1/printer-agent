package main

import (
	"os"
	"strings"

	"printer-agent/server"
)

func main() {
	if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "aharsuchi-printer://") {
		server.HandleProtocol(os.Args[1])
		return
	}

	// Optional: still allow HTTP mode
	// server.Start()
}
