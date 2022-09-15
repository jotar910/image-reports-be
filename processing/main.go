package main

import (
	"fmt"
	"image-reports/processing/configs"
	"image-reports/processing/pkg/transport"
	"os"

	log "image-reports/helpers/services/logger"
	"image-reports/helpers/services/server"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	// Connect to DB
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=GMT",
		config.ServiceConfig.Database.Host,
		config.ServiceConfig.Database.Username,
		config.ServiceConfig.Database.Password,
		config.ServiceConfig.Database.Database,
		config.ServiceConfig.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
		os.Exit(1)
	}

	// Create and run server
	server.NewServer(transport.NewServerConfiguration(db, config)).Run()
}
