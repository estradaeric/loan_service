package models

type Investor struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`
}