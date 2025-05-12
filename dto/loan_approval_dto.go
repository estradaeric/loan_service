package dto

import "mime/multipart"

type LoanApprovalDTO struct {
	FieldValidatorID string                  `form:"field_validator_id" validate:"required"`
	ProofImage       *multipart.FileHeader   `form:"proof_image" validate:"required"`
}