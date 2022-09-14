package main

import (
	"image-reports/processing/configs"
	"image-reports/processing/pkg/transport"

	log "image-reports/helpers/services/logger"
	"image-reports/helpers/services/server"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Logger
	log.Initialize()

	// Initialize Configs
	config, err := configs.Initialize("processing")
	if err != nil {
		log.Fatalf("config: %s", err)
	}
	log.Infof("Starting with config: %+v", config)

	gin.SetMode(config.Gin.Mode)

	// Create and run server
	server.NewServer(transport.NewServerConfiguration(config)).Run()
}
