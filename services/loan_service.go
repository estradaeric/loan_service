package services

import (
	"errors"
	"fmt"
	"time"

	"loan-service/config"
	"loan-service/dto"
	"loan-service/logger"
	"loan-service/models"
	"loan-service/repository"
	"loan-service/utils"

	"github.com/sirupsen/logrus"

)

// LoanService defines business logic for handling loans
type LoanService struct {
	repo repository.LoanRepository
	cfg  *config.Config
}

// NewLoanService creates a new LoanService with the given repository
func NewLoanService(repo repository.LoanRepository, cfg *config.Config) *LoanService {
	return &LoanService{
		repo: repo, 
		cfg: cfg,
	}
}

func (s *LoanService) GetLoanByID(id string) (*models.Loan, error) {
	return s.repo.GetByID(id)
}

// CreateLoan initializes a loan creation
func (s *LoanService) CreateLoan(input *dto.LoanCreateDTO) (*models.Loan, error) {
	loan := &models.Loan{
		ID:         utils.GenerateID("loan_"),
		State:      models.LoanStateProposed,
		Principal:  input.Principal,
		Rate:       input.Rate,
		ROI:        input.ROI,
		BorrowerID: input.BorrowerID,
		CreatedAt:  time.Now(),
	}

	if err := loan.ValidateInitialState(); err != nil {
		logger.Log.WithError(err).Warn("Loan validation failed")
		return nil, err
	}

	if err := s.repo.Save(loan); err != nil {
		logger.Log.WithError(err).Error("Failed to save loan")
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"loan_id": loan.ID,
		"state":   loan.State,
	}).Info("Loan created successfully")

	return loan, nil
}

// ApproveLoan moves a loan to 'approved' state
func (s *LoanService) ApproveLoan(loanID string, input *dto.LoanApprovalDTO, proof []byte) (*models.Loan, error) {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		logger.Log.WithError(err).Error("Loan not found")
		return nil, err
	}
	if loan.State != models.LoanStateProposed {
		return nil, errors.New("loan must be in 'proposed' state to be approved")
	}
	if loan.ApprovedAt != nil {
		return nil, errors.New("loan already approved")
	}

	now := time.Now()
	loan.State = models.LoanStateApproved
	loan.ApprovedAt = &now
	loan.FieldValidatorID = &input.FieldValidatorID
	loan.ProofImage = proof

	if err := s.repo.Update(loan); err != nil {
		return nil, err
	}

	logger.Log.WithField("loan_id", loan.ID).Info("Loan approved")
	return loan, nil
}

// InvestLoan handles investment and over-invest validation
func (s *LoanService) InvestLoan(loanID string, input *dto.LoanInvestDTO) (*models.Loan, error) {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get loan for investment")
		return nil, err
	}

	if loan.State != models.LoanStateApproved {
		logger.Log.Warnf("Loan %s is not in 'approved' state", loan.ID)
		return nil, errors.New("loan must be in 'approved' state")
	}

	if input.Amount <= 0 {
		logger.Log.Warn("Investment amount must be greater than 0")
		return nil, errors.New("invalid amount")
	}

	for _, inv := range loan.Investments {
		if inv.InvestorID == input.InvestorID {
			logger.Log.Warnf("Investor %s has already invested in loan %s", input.InvestorID, loan.ID)
			return nil, errors.New("investor already invested")
		}
	}

	total := loan.TotalInvested()
	if total+input.Amount > loan.Principal {
		logger.Log.Warnf("Investment exceeds principal: %.2f + %.2f > %.2f", total, input.Amount, loan.Principal)
		return nil, errors.New("investment exceeds loan principal")
	}

	investor, err := s.repo.GetInvestorByID(input.InvestorID)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to get investor data")
		return nil, err
	}

	// Generate agreement letter as PDF
	pdfContent, err := utils.GenerateAgreementPDF(loan.ID, investor.Name, input.Amount)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to generate agreement PDF")
		return nil, err
	}

	// Assign to model for DB storage
	loan.AgreementLetter = pdfContent

	investment := models.LoanInvestment{
		ID:         utils.GenerateID("inv_"),
		LoanID:     loan.ID,
		InvestorID: investor.ID,
		Amount:     input.Amount,
		CreatedAt:  time.Now(),
		Investor:   investor,
	}
	loan.Investments = append(loan.Investments, investment)

	if loan.TotalInvested() == loan.Principal {
		loan.State = models.LoanStateInvested
		logger.Log.Infof("Loan %s fully funded", loan.ID)
	}

	if err := s.repo.Update(loan); err != nil {
		logger.Log.WithError(err).Error("Failed to update loan after investment")
		return nil, err
	}

	// Send email with PDF attachment
	emailBody := fmt.Sprintf("Hi %s,\n\nThank you for your investment in loan %s.\nPlease find the agreement letter attached.\n\nRegards,\nLoan Service", investor.Name, loan.ID)
	if err := utils.SendEmailWithAttachment(investor.Email, "Loan Agreement", emailBody, pdfContent, "agreement.pdf"); err != nil {
		logger.Log.WithError(err).Error("Failed to send agreement email to investor")
	}

	logger.Log.WithFields(logrus.Fields{
		"loan_id":     loan.ID,
		"investor_id": investor.ID,
		"amount":      input.Amount,
	}).Info("Investment processed successfully")

	return loan, nil
}

// DisburseLoan finalizes loan disbursement
func (s *LoanService) DisburseLoan(loanID string, input *dto.LoanDisburseDTO, signed []byte) (*models.Loan, error) {
	loan, err := s.repo.GetByID(loanID)
	if err != nil {
		return nil, err
	}
	if loan.State != models.LoanStateInvested {
		return nil, errors.New("loan must be in 'invested' state")
	}
	if loan.DisbursedAt != nil {
		return nil, errors.New("loan already disbursed")
	}

	officer, err := s.repo.GetEmployeeByID(input.FieldOfficerID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	loan.State = models.LoanStateDisbursed
	loan.DisbursedAt = &now
	loan.DisburseNotes = input.DisburseNotes
	loan.SignedAgreement = signed
	loan.FieldOfficerID = &input.FieldOfficerID
	loan.FieldOfficer = officer

	if err := s.repo.Update(loan); err != nil {
		return nil, err
	}

	logger.Log.WithField("loan_id", loan.ID).Info("Loan disbursed")
	return loan, nil
}