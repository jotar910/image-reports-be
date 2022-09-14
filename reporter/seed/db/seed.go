package main

import (
	"fmt"
	"os"

	"image-reports/reporter/configs"
	"image-reports/reporter/models"

	log "image-reports/helpers/services/logger"

	shared_models "image-reports/shared/models"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize Logger
	log.Initialize()

	// Initialize Configs
	config, err := configs.Initialize("reporter")
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
	tables := []any{&models.Reports{}, &models.Approvals{}}
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

	// Add reports to database.
	reports := []models.Reports{
		{
			Name:    "Report #1",
			UserID:  1,
			ImageID: "image-1-uuid",
			Status:  shared_models.ReportStatusPublished,
			Approval: models.Approvals{
				UserID: 1,
				Status: shared_models.ApprovalStatusApproval,
			},
		},
		{
			Name:    "Report #2",
			UserID:  1,
			ImageID: "image-2-uuid",
			Status:  shared_models.ReportStatusPublished,
			Approval: models.Approvals{
				UserID: 1,
				Status: shared_models.ApprovalStatusRejected,
			},
		},
		{
			Name:    "Report #3",
			UserID:  2,
			ImageID: "image-3-uuid",
			Status:  shared_models.ReportStatusNew,
		},
		{
			Name:    "Report #4",
			UserID:  3,
			ImageID: "image-4-uuid",
			Status:  shared_models.ReportStatusEvaluating,
		},
		{
			Name:    "Report #5",
			UserID:  1,
			ImageID: "image-5-uuid",
			Status:  shared_models.ReportStatusError,
		},
		{
			Name:    "Report #6",
			UserID:  2,
			ImageID: "image-6-uuid",
			Status:  shared_models.ReportStatusPending,
		},
		{
			Name:    "Report #7",
			UserID:  3,
			ImageID: "image-7-uuid",
			Status:  shared_models.ReportStatusPending,
		},
	}
	tx := db.Create(reports)
	if tx.Error != nil {
		log.Fatalf("failed to create reports: %s", err.Error)
		os.Exit(1)
	}
	log.Infof("Created report records: %+v", reports)

	os.Exit(0)
}
