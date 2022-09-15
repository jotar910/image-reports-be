package main

import (
	"fmt"
	"os"

	"image-reports/processing/configs"
	"image-reports/processing/models"

	log "image-reports/helpers/services/logger"

	"github.com/spf13/cobra"
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

	// Connect to DB
	log.Info("Connecting to database...")
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

	// Create command
	command := &cobra.Command{}
	flags := command.Flags()
	flags.BoolP("auto-migrate", "a", false, "auto migrate the database")
	flags.BoolP("keep-state", "k", false, "keep database old state")

	migrator := db.Migrator()
	tables := []any{&models.Evaluations{}, &models.Categories{}, &models.EvaluationCategories{}}
	log.Info("Running GORM seed...")

	// Auto migrate.
	if auto, err := flags.GetBool("auto-migrate"); err != nil {
		panic(err)
	} else if auto {
		if err := migrator.AutoMigrate(tables...); err != nil {
			log.Fatalf("failed to run auto migration: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Keep old state.
	if keepState, err := flags.GetBool("keep-state"); err != nil {
		panic(err)
	} else if !keepState {
		if err := migrator.DropTable(tables...); err != nil {
			log.Fatalf("failed to drop tables: %v", err)
			os.Exit(1)
		}

		if err := migrator.CreateTable(tables...); err != nil {
			log.Fatalf("failed to create new tables: %v", err)
			os.Exit(1)
		}
	}

	log.Info("Adding records...")

	// Add categories to database.
	categories := []models.Categories{
		{Name: "Cinematic"},
		{Name: "Mythic"},
		{Name: "Violence"},
		{Name: "Abuse"},
	}
	tx := db.Create(categories)
	if tx.Error != nil {
		log.Fatalf("failed to create categories: %s", err.Error)
		os.Exit(1)
	}
	log.Infof("Created category records: %+v", categories)

	// Add evaluations to database.
	evaluations := []models.Evaluations{
		{
			ReportID:   1,
			ImageID:    "testing",
			Grade:      75,
			Categories: []models.Categories{categories[0], categories[1]},
		},
		{
			ReportID:   2,
			ImageID:    "testing",
			Grade:      20,
			Categories: []models.Categories{categories[1], categories[2], categories[3]},
		},
	}
	tx = db.Create(evaluations)
	if tx.Error != nil {
		log.Fatalf("failed to create evaluations: %s", err.Error)
		os.Exit(1)
	}
	log.Infof("Created evaluation records: %+v", evaluations)

	os.Exit(0)
}
