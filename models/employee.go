package models

// Employee represents a staff member who acts as a field validator
type Employee struct {
	ID   string `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Email string `json:"email" gorm:"uniqueIndex"`	
}