package main

import (
	"image-reports/helpers/services/server"
	"image-reports/realtime/pkg/transport"
)

func main() {
	// Create and run server
	server.NewServer(transport.NewServerConfiguration()).Run()
}
