package main

import (
	"fmt"
	"os"

	"image-reports/users/configs"
	"image-reports/users/models"

	"image-reports/helpers/services/auth"
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
	config, err := configs.Initialize("users")
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
	tables := []any{&models.Roles{}, &models.Users{}}
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

	// Add users to database.
	roles := []models.Roles{
		{Name: shared_models.AdminRole},
		{Name: shared_models.UserRole},
	}
	tx := db.Create(roles)
	if tx.Error != nil {
		log.Fatalf("failed to create roles: %s", err.Error)
		os.Exit(1)
	}
	log.Infof("Created role records: %+v", roles)

	// Add users to database.
	// TODO: generate random passwords.
	users := []models.Users{
		{
			Email:    "admin@email.com",
			Password: "admin",
			Role:     roles[0],
		},
		{
			Email:    "user@email.com",
			Password: "user",
			Role:     roles[1],
		},
	}
	for i, user := range users {
		encrypted, err := auth.HashPassword(user.Password)
		if err != nil {
			log.Fatalf("failed to encrypt password for user %s: %v", user.Email, err)
		}
		users[i].Password = encrypted
	}
	tx = db.Create(users)
	if tx.Error != nil {
		log.Fatalf("failed to create users: %s", err.Error)
		os.Exit(1)
	}
	log.Infof("Created user records: %+v", users)

	os.Exit(0)
}
