package main

import (
	"image-reports/users/pkg/transport"

	"image-reports/helpers/services/server"
)

func main() {
	// Create and run server
	server.NewServer(transport.NewServerConfiguration()).Run()
}
