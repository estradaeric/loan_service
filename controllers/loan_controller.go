package controllers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"loan-service/dto"
	"loan-service/services"
	"loan-service/utils"
)

// LoanController handles HTTP requests related to loans
type LoanController struct {
	Service *services.LoanService
}

// Validator instance
var validate = validator.New()

// NewLoanController initializes a new LoanController
func NewLoanController(svc *services.LoanService) *LoanController {  
	return &LoanController{Service: svc}
}

// CreateLoan handles POST /loans
func (lc *LoanController) CreateLoan(w http.ResponseWriter, r *http.Request) {
	var input dto.LoanCreateDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	loan, err := lc.Service.CreateLoan(&input)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteSuccess(w, http.StatusCreated, "Loan created", dto.ToLoanResponseDTO(loan))
}

// ApproveLoan handles PUT /loans/{id}/approve with file upload (proof image)
func (lc *LoanController) ApproveLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	// Validate fields
	input := dto.LoanApprovalDTO{
		FieldValidatorID: r.FormValue("field_validator_id"),
	}
	file, fileHeader, err := r.FormFile("proof_image")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Proof image is required")
		return
	}
	defer file.Close()

	if err := utils.IsValidImage(fileHeader); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	imageData, err := io.ReadAll(file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to read image")
		return
	}

	loan, err := lc.Service.ApproveLoan(id, &input, imageData)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Loan approved", dto.ToLoanResponseDTO(loan))
}

// InvestLoan handles POST /loans/{id}/invest
func (lc *LoanController) InvestLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	loanID := vars["id"]

	var input dto.LoanInvestDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := validate.Struct(&input); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	loan, err := lc.Service.InvestLoan(loanID, &input)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.WriteSuccess(w, http.StatusOK, "Investment added", dto.ToLoanResponseDTO(loan))
}

// DisburseLoan handles PUT /loans/{id}/disburse with file upload (signed agreement)
func (lc *LoanController) DisburseLoan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to parse form: "+err.Error())
		return
	}

	input := dto.LoanDisburseDTO{
		FieldOfficerID: r.FormValue("field_officer_id"),
		DisburseNotes:  r.FormValue("disburse_notes"),
	}

	file, fileHeader, err := r.FormFile("signed_agreement")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Signed agreement is required")
		return
	}
	defer file.Close()

	if err := utils.IsValidImage(fileHeader); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	agreementData, err := io.ReadAll(file)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to read file")
		return
	}

	loan, err := lc.Service.DisburseLoan(id, &input, agreementData)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, "Loan disbursed", dto.ToLoanResponseDTO(loan))
}

// Get Agreement File
func (lc *LoanController) GetAgreement(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	loan, err := lc.Service.GetLoanByID(id)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Loan not found")
		return
	}

	if loan.AgreementLetter == nil {
		utils.WriteError(w, http.StatusNotFound, "No file found")
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=agreement.pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(loan.AgreementLetter)
}