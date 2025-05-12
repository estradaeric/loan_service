package models

// Borrower represents the person receiving the loan
type Borrower struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	NIK  string `json:"nik" gorm:"uniqueIndex"`
}