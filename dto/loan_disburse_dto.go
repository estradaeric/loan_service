package dto

type LoanDisburseDTO struct {
	FieldOfficerID string `json:"field_officer_id" validate:"required"`
	DisburseNotes  string `json:"disburse_notes"`
}