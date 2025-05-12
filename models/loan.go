package models

import (
	"errors"
	"time"
)

type LoanState string

const (
	LoanStateProposed  LoanState = "proposed"
	LoanStateApproved  LoanState = "approved"
	LoanStateInvested  LoanState = "invested"
	LoanStateDisbursed LoanState = "disbursed"
)

type Loan struct {
	ID         			string				`json:"id" gorm:"primaryKey"`
	State      			LoanState			`json:"state"`
	Principal  			float64				`json:"principal"`
	Rate       			float64				`json:"rate"`
	ROI        			float64				`json:"roi"`
	CreatedAt  			time.Time			`json:"created_at"`
	ApprovedAt 			*time.Time			`json:"approved_at"`
	DisbursedAt 		*time.Time			`json:"disbursed_at"`

	// Foreign Keys
	BorrowerID			string				`json:"borrower_id"`
	Borrower			*Borrower			`json:"-" gorm:"foreignKey:BorrowerID"`

	FieldValidatorID	*string				`json:"field_validator_id,omitempty"`
	FieldValidator   	*Employee			`json:"-" gorm:"foreignKey:FieldValidatorID"`

	FieldOfficerID		*string				`json:"field_officer_id,omitempty"`
	FieldOfficer		*Employee			`json:"-" gorm:"foreignKey:FieldOfficerID"`

	ProofImage      	[]byte 				`json:"-" gorm:"type:bytea"`
	AgreementLetter 	[]byte 				`json:"-" gorm:"type:bytea"`
	SignedAgreement 	[]byte 				`json:"-" gorm:"type:bytea"`
	DisburseNotes		string     			`json:"disburse_notes,omitempty"`
	
	// Relationship One-to-many
	Investments 		[]LoanInvestment	`json:"investments,omitempty" gorm:"foreignKey:LoanID"`

}

// IsValidState checks if the loan state is valid
func IsValidState(state LoanState) bool {
	switch state {
	case LoanStateProposed, LoanStateApproved, LoanStateInvested, LoanStateDisbursed:
		return true
	default:
		return false
	}
}

// ValidateInitialState ensures loan is correctly initialized
func (l *Loan) ValidateInitialState() error {
	if l.BorrowerID == "" {
		return errors.New("borrower_id is required")
	}
	if l.Principal <= 0 {
		return errors.New("principal must be greater than 0")
	}
	if !IsValidState(l.State) {
		return errors.New("invalid loan state")
	}
	if l.State != LoanStateProposed {
		return errors.New("loan must start in 'proposed' state")
	}
	return nil
}

// TotalInvested calculates the total amount invested in the loan
func (l *Loan) TotalInvested() float64 {
	var total float64
	for _, inv := range l.Investments {
		total += inv.Amount
	}
	return total
}