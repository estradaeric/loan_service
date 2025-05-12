package repository

import "loan-service/models"

// LoanRepository defines the behavior required for storing and retrieving loans.
// This interface is implemented by different storage backends (e.g., memory or PostgreSQL).
type LoanRepository interface {
	Save(loan *models.Loan) error
	GetByID(id string) (*models.Loan, error)
	Update(loan *models.Loan) error
	
    GetInvestorByID(id string) (*models.Investor, error)
    GetEmployeeByID(id string) (*models.Employee, error)	
}