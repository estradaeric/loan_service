package repository

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"loan-service/models"
)

type PostgresLoanRepository struct {
	db *gorm.DB
}

func NewPostgresLoanRepository(dsn string) (*PostgresLoanRepository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.Borrower{},
		&models.Employee{},
		&models.Investor{},
		&models.Loan{},
		&models.LoanInvestment{},
	); err != nil {
		return nil, err
	}

	return &PostgresLoanRepository{db: db}, nil
}

func (r *PostgresLoanRepository) Save(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

func (r *PostgresLoanRepository) Update(loan *models.Loan) error {
	return r.db.Save(loan).Error
}

func (r *PostgresLoanRepository) GetByID(id string) (*models.Loan, error) {
	var loan models.Loan
	err := r.db.
		Preload("Borrower").
		Preload("FieldValidator").
		Preload("FieldOfficer").
		Preload("Investments").
		Preload("Investments.Investor").
		First(&loan, "id = ?", id).Error

	if err != nil {
		return nil, errors.New("loan not found")
	}
	return &loan, nil
}

func (r *PostgresLoanRepository) GetInvestorByID(id string) (*models.Investor, error) {
	var investor models.Investor
	if err := r.db.First(&investor, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &investor, nil
}

func (r *PostgresLoanRepository) GetEmployeeByID(id string) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.First(&employee, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

func NewPostgresLoanRepositoryFromDB(db *gorm.DB) *PostgresLoanRepository {
	return &PostgresLoanRepository{db: db}
}