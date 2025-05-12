package dto

type LoanInvestDTO struct {
	InvestorID string  `json:"investor_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required,gt=0"`
}