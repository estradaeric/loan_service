package repository

import (
	"errors"
	"sync"

	"loan-service/models"
)

// MemoryRepository holds all in-memory data stores and their logic
type MemoryRepository struct {
	mu              sync.RWMutex
	loanStore       map[string]*models.Loan
	investStore     map[string][]models.LoanInvestment
	investorStore   map[string]*models.Investor
	borrowerStore   map[string]*models.Borrower
	employeeStore   map[string]*models.Employee
}

// NewMemoryRepository initializes the in-memory repository
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		loanStore:     make(map[string]*models.Loan),
		investStore:   make(map[string][]models.LoanInvestment),
		investorStore: make(map[string]*models.Investor),
		borrowerStore: make(map[string]*models.Borrower),
		employeeStore: make(map[string]*models.Employee),
	}
}

//
// ========== LOAN SECTION ==========
//

// SaveLoan saves a new loan
func (r *MemoryRepository) SaveLoan(loan *models.Loan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.loanStore[loan.ID]; exists {
		return errors.New("loan already exists")
	}
	r.loanStore[loan.ID] = loan

	if len(loan.Investments) > 0 {
		r.investStore[loan.ID] = loan.Investments
	}

	return nil
}

// GetLoanByID fetches loan and related investments + investor
func (r *MemoryRepository) GetLoanByID(id string) (*models.Loan, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	loan, exists := r.loanStore[id]
	if !exists {
		return nil, errors.New("loan not found")
	}

	clone := *loan
	clone.Investments = nil

	if investments, ok := r.investStore[id]; ok {
		for _, inv := range investments {
			if investor, found := r.investorStore[inv.InvestorID]; found {
				inv.Investor = investor
			}
			clone.Investments = append(clone.Investments, inv)
		}
	}

	return &clone, nil
}

// UpdateLoan updates an existing loan
func (r *MemoryRepository) UpdateLoan(loan *models.Loan) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.loanStore[loan.ID]; !exists {
		return errors.New("loan not found")
	}

	r.loanStore[loan.ID] = loan

	if loan.Investments != nil {
		r.investStore[loan.ID] = loan.Investments
	}

	return nil
}

//
// ========== INVESTOR SECTION ==========
//

// SaveInvestor saves a new investor
func (r *MemoryRepository) SaveInvestor(inv *models.Investor) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.investorStore[inv.ID]; exists {
		return errors.New("investor already exists")
	}
	r.investorStore[inv.ID] = inv
	return nil
}

// GetInvestorByID fetches an investor
func (r *MemoryRepository) GetInvestorByID(id string) (*models.Investor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	inv, exists := r.investorStore[id]
	if !exists {
		return nil, errors.New("investor not found")
	}
	return inv, nil
}

//
// ========== BORROWER SECTION ==========
//

// SaveBorrower saves a new borrower
func (r *MemoryRepository) SaveBorrower(b *models.Borrower) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.borrowerStore[b.ID]; exists {
		return errors.New("borrower already exists")
	}
	r.borrowerStore[b.ID] = b
	return nil
}

// GetBorrowerByID fetches a borrower
func (r *MemoryRepository) GetBorrowerByID(id string) (*models.Borrower, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	b, exists := r.borrowerStore[id]
	if !exists {
		return nil, errors.New("borrower not found")
	}
	return b, nil
}

//
// ========== EMPLOYEE SECTION ==========
//

// SaveEmployee saves a new employee
func (r *MemoryRepository) SaveEmployee(e *models.Employee) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.employeeStore[e.ID]; exists {
		return errors.New("employee already exists")
	}
	r.employeeStore[e.ID] = e
	return nil
}

// GetEmployeeByID fetches an employee
func (r *MemoryRepository) GetEmployeeByID(id string) (*models.Employee, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	e, exists := r.employeeStore[id]
	if !exists {
		return nil, errors.New("employee not found")
	}
	return e, nil
}

func (r *MemoryRepository) Save(loan *models.Loan) error {
	return r.SaveLoan(loan)
}

func (r *MemoryRepository) GetByID(id string) (*models.Loan, error) {
	return r.GetLoanByID(id)
}

func (r *MemoryRepository) Update(loan *models.Loan) error {
	return r.UpdateLoan(loan)
}