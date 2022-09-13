package main

import (
	"image-reports/helpers/services/server"
	"image-reports/reporter/pkg/transport"
)

func main() {
	// Create and run server
	server.NewServer(transport.NewServerConfiguration()).Run()
}
