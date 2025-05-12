package main

import (
	"fmt"
	"net/http"
	"os"

	"loan-service/config"
	"loan-service/controllers"
	"loan-service/logger"
	"loan-service/repository"
	"loan-service/routes"
	"loan-service/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.InitLogger("development") // Fallback to development if config error
		logger.Log.WithError(err).Fatal("Failed to load config")
		os.Exit(1)
	}

	// Initialize logger based on environment
	logger.InitLogger(cfg.Env)
	logger.Log.Infof("Starting Loan Service in %s environment", cfg.Env)

	// Prepare DSN string for DB connection
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)
	logger.Log.Debugf("Using DSN: %s", dsn)

	var loanRepo repository.LoanRepository

	// Attempt PostgreSQL connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Log.WithError(err).Warn("PostgreSQL connection failed, falling back to in-memory repository")

		memRepo := repository.NewMemoryRepository()
		memRepo.SeedMemoryData()
		loanRepo = memRepo
	} else {
		logger.Log.Info("Connected to PostgreSQL successfully")

		if err := repository.AutoMigrateTables(db); err != nil {
			logger.Log.WithError(err).Fatal("Failed to migrate database tables")
		}

		if err := repository.SeedData(db); err != nil {
			logger.Log.WithError(err).Fatal("Failed to seed initial data")
		}

		pgRepo := repository.NewPostgresLoanRepositoryFromDB(db)
		loanRepo = pgRepo
	}

	// Initialize services, controllers, and router
	service := services.NewLoanService(loanRepo, cfg)
	controller := controllers.NewLoanController(service)
	router := routes.SetupRouter(controller)

	// Start HTTP server
	logger.Log.Infof("Loan service started on port %s in %s environment", cfg.Port, cfg.Env)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		logger.Log.WithError(err).Fatal("Failed to start server")
	}
}