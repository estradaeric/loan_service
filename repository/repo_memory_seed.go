package repository

import "loan-service/models"

// SeedMemoryData seeds default dummy data into memory store
func (r *MemoryRepository) SeedMemoryData() {
	// Borrowers
	borrowers := []models.Borrower{
		{ID: "BORR-001", Name: "Agus Suryanto", NIK: "3271000001"},
		{ID: "BORR-002", Name: "Dewi Kartika", NIK: "3271000002"},
		{ID: "BORR-003", Name: "Budi Gunawan", NIK: "3271000003"},
		{ID: "BORR-004", Name: "Fitriani Sari", NIK: "3271000004"},
		{ID: "BORR-005", Name: "Rizki Maulana", NIK: "3271000005"},
		{ID: "BORR-006", Name: "Intan Nuraini", NIK: "3271000006"},
		{ID: "BORR-007", Name: "Sigit Prabowo", NIK: "3271000007"},
		{ID: "BORR-008", Name: "Yuni Astuti", NIK: "3271000008"},
		{ID: "BORR-009", Name: "Teguh Santosa", NIK: "3271000009"},
		{ID: "BORR-010", Name: "Lilis Rahayu", NIK: "3271000010"},
	}
	for _, b := range borrowers {
		r.SaveBorrower(&b)
	}

	// Employees
	employees := []models.Employee{
		{ID: "EMP-001", Name: "Maya Kusuma", Email: "maya.kusuma@company.com"},
		{ID: "EMP-002", Name: "Dian Prasetyo", Email: "dian.prasetyo@company.com"},
		{ID: "EMP-003", Name: "Rani Febriani", Email: "rani.febriani@company.com"},
		{ID: "EMP-004", Name: "Fajar Rahman", Email: "fajar.rahman@company.com"},
		{ID: "EMP-005", Name: "Nina Sari", Email: "nina.sari@company.com"},
	}
	for _, e := range employees {
		r.SaveEmployee(&e)
	}

	// Investors
	investors := []models.Investor{
		{ID: "INV-001", Name: "Randy Saputra", Email: "randy@example.com"},
		{ID: "INV-002", Name: "Wulan Kartika", Email: "wulan@example.com"},
		{ID: "INV-003", Name: "Andika Mahesa", Email: "andika@example.com"},
		{ID: "INV-004", Name: "Nadya Ayu", Email: "nadya@example.com"},
		{ID: "INV-005", Name: "Ilham Ramadhan", Email: "ilham@example.com"},
	}
	for _, inv := range investors {
		r.SaveInvestor(&inv)
	}
}
