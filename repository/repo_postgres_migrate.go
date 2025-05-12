package repository

import (
	"loan-service/models"

	"gorm.io/gorm"
)

// AutoMigrateTables performs schema migration for all models in the correct order
func AutoMigrateTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Borrower{},
		&models.Employee{},
		&models.Investor{},
		&models.Loan{},
		&models.LoanInvestment{},
	)
}