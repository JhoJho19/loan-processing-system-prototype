package service

import "time"

type CreateApplicationInput struct {
	FirstName       string
	LastName        string
	BirthDate       time.Time
	MonthlyIncome   float64
	EmploymentType  string
	Email           string
	Phone           string
	Address         string
	RequestedAmount int
	TermMonths      int
}
