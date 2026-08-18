package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidApplicationID = errors.New("invalid application ID")
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
