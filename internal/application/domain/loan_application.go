package domain

import (
	"time"
)

type LoanApplication struct {
	ID              string
	CustomerID      string
	RequestedAmount int
	TermMonths      int
	Status          ApplicationStatus
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
