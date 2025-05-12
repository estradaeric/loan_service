package utils

import (
	"bytes"
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

func GenerateAgreementPDF(loanID, investorName string, amount float64) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Loan Agreement")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(190, 10, fmt.Sprintf("This agreement certifies that investor %s has committed %.2f to loan %s.", investorName, amount, loanID), "", "", false)

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}