package domain

import (
	"time"
)

type Customer struct {
	ID             string
	FirstName      string
	LastName       string
	BirthDate      time.Time
	MonthlyIncome  int
	EmploymentType string
	Email          string
	Phone          string
	Address        string
}
