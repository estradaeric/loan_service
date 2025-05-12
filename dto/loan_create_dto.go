package dto

type LoanCreateDTO struct {
	BorrowerID string  `json:"borrower_id" validate:"required"`
	Principal  float64 `json:"principal" validate:"required,gt=0"`
	Rate       float64 `json:"rate" validate:"required,gt=0"`
	ROI        float64 `json:"roi" validate:"required,gt=0"`
}