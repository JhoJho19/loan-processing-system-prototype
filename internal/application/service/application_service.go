package service

import (
	"loan-processing-system/internal/application/domain"
	"time"

	"github.com/google/uuid"
)

type ApplicationService struct{}

func (c ApplicationService) CreateApplication(customerApplicationInput CreateApplicationInput) (domain.LoanApplication, error) {

	customer := domain.Customer{}
	customer.ID = uuid.NewString()
	customer.FirstName = customerApplicationInput.FirstName
	customer.LastName = customerApplicationInput.LastName
	customer.Address = customerApplicationInput.Address
	customer.BirthDate = customerApplicationInput.BirthDate
	customer.MonthlyIncome = customerApplicationInput.MonthlyIncome
	customer.EmploymentType = customerApplicationInput.EmploymentType
	customer.Email = customerApplicationInput.Email
	customer.Phone = customerApplicationInput.Phone

	loanApplication := domain.LoanApplication{}
	loanApplication.ID = uuid.NewString()
	loanApplication.CustomerID = customer.ID
	loanApplication.RequestedAmount = customerApplicationInput.RequestedAmount
	loanApplication.TermMonths = customerApplicationInput.TermMonths
	loanApplication.Status = domain.StatusNew

	timeNow := time.Now()
	loanApplication.CreatedAt = timeNow
	loanApplication.UpdatedAt = timeNow

	return loanApplication, nil
}

func (c ApplicationService) GetApplication(applicationID string) (domain.LoanApplication, error) {

	return domain.LoanApplication{}, nil
}
