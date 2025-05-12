package dto

import (
	"time"
	"loan-service/models"
)

type LoanResponseDTO struct {
	ID               string     `json:"id"`
	State            string     `json:"state"`
	Principal        float64    `json:"principal"`
	Rate             float64    `json:"rate"`
	ROI              float64    `json:"roi"`
	InterestAmount   float64    `json:"interest_amount"`
	CreatedAt        time.Time  `json:"created_at"`
	ApprovedAt       *time.Time `json:"approved_at,omitempty"`
	DisbursedAt      *time.Time `json:"disbursed_at,omitempty"`
	FieldValidatorID string     `json:"field_validator_id,omitempty"`
	FieldOfficerID   string     `json:"field_officer_id,omitempty"`
	ProofUploaded    bool       `json:"proof_uploaded"`
	AgreementExists  bool       `json:"agreement_exists"`
	SignedUploaded   bool       `json:"signed_uploaded"`
	DisburseNotes    string     `json:"disburse_notes,omitempty"`

	Borrower    *BorrowerSimpleDTO  `json:"borrower,omitempty"`
	Investments []LoanInvestmentDTO `json:"investments,omitempty"`
}

type BorrowerSimpleDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LoanInvestmentDTO struct {
	ID            string    `json:"id"`
	InvestorID    string    `json:"investor_id"`
	InvestorName  string    `json:"investor_name,omitempty"`
	InvestorEmail string    `json:"investor_email,omitempty"`
	Amount        float64   `json:"amount"`
	InvestorProfit float64  `json:"investor_profit"`
	CreatedAt     time.Time `json:"created_at"`
}

func ToLoanResponseDTO(l *models.Loan) *LoanResponseDTO {
	dto := &LoanResponseDTO{
		ID:               l.ID,
		State:            string(l.State),
		Principal:        l.Principal,
		Rate:             l.Rate,
		ROI:              l.ROI,
		InterestAmount:   l.Principal * l.Rate / 100,
		CreatedAt:        l.CreatedAt,
		ApprovedAt:       l.ApprovedAt,
		DisbursedAt:      l.DisbursedAt,
		FieldValidatorID: deref(l.FieldValidatorID),
		FieldOfficerID:   deref(l.FieldOfficerID),
		ProofUploaded:    len(l.ProofImage) > 0,
		AgreementExists:  len(l.AgreementLetter) > 0,
		SignedUploaded:   len(l.SignedAgreement) > 0,
		DisburseNotes:    l.DisburseNotes,
	}

	if l.Borrower != nil {
		dto.Borrower = &BorrowerSimpleDTO{
			ID:   l.Borrower.ID,
			Name: l.Borrower.Name,
		}
	}

	for _, inv := range l.Investments {
		invDTO := LoanInvestmentDTO{
			ID:         inv.ID,
			InvestorID: inv.InvestorID,
			Amount:     inv.Amount,
			InvestorProfit: inv.Amount * l.ROI / 100,
			CreatedAt:  inv.CreatedAt,
		}

		if inv.Investor != nil {
			invDTO.InvestorName = inv.Investor.Name
			invDTO.InvestorEmail = inv.Investor.Email
		}

		dto.Investments = append(dto.Investments, invDTO)
	}

	return dto
}

func deref(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}