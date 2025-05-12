package models

import "time"

type LoanInvestment struct {
	ID         string     `json:"id" gorm:"primaryKey"`
	LoanID     string     `json:"loan_id"`
	InvestorID string     `json:"investor_id"`
	Amount     float64    `json:"amount"`
	CreatedAt  time.Time  `json:"created_at"`

	// Relations
	Loan     *Loan     `json:"-" gorm:"foreignKey:LoanID"`
	Investor *Investor `json:"-" gorm:"foreignKey:InvestorID"`

}