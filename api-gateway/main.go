package main

import (
	"image-reports/api-gateway/configs"
	"image-reports/api-gateway/pkg/transport"

	log "image-reports/helpers/services/logger"
	"image-reports/helpers/services/server"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Logger
	log.Initialize()

	// Initialize Configs
	config, err := configs.Initialize("users")
	if err != nil {
		log.Fatalf("config: %s", err)
	}
	log.Infof("Starting with config: %+v", config)

	gin.SetMode(config.Gin.Mode)

	// Create and run server
	server.NewServer(transport.NewServerConfiguration(config)).Run()
}
