package service

import (
	"errors"
	"loan-processing-system/internal/application/domain"
	"loan-processing-system/internal/application/repository"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ApplicationService struct {
	repository repository.ApplicationRepository
}

func NewApplicationService(repo repository.ApplicationRepository) ApplicationService {
	return ApplicationService{
		repository: repo,
	}
}

func (s ApplicationService) CreateApplication(customerApplicationInput CreateApplicationInput) (domain.LoanApplication, error) {

	if err := validateCreateApplicationInput(customerApplicationInput); err != nil {
		return domain.LoanApplication{}, err
	}

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

	repositoryErr := s.repository.Create(customer, loanApplication)
	if repositoryErr != nil {
		return domain.LoanApplication{}, repositoryErr
	}

	return loanApplication, nil
}

func (s ApplicationService) GetApplication(applicationID string) (domain.LoanApplication, error) {

	if applicationID == "" {
		return domain.LoanApplication{}, domain.ErrInvalidApplicationID
	}

	loanApplication, repositoryErr := s.repository.GetByID(applicationID)
	if repositoryErr != nil {
		return domain.LoanApplication{}, repositoryErr
	}

	return loanApplication, nil
}

func validateCreateApplicationInput(input CreateApplicationInput) error {
	errorsSlice := []string{}

	if input.FirstName == "" {
		errorsSlice = append(errorsSlice, "first name is required")
	}
	if input.LastName == "" {
		errorsSlice = append(errorsSlice, "last name is required")
	}
	if input.Address == "" {
		errorsSlice = append(errorsSlice, "address is required")
	}
	if input.BirthDate.IsZero() {
		errorsSlice = append(errorsSlice, "birth date is required")
	}
	if input.MonthlyIncome <= 0 {
		errorsSlice = append(errorsSlice, "monthly income must be greater than zero")
	}
	if input.EmploymentType == "" {
		errorsSlice = append(errorsSlice, "employment type is required")
	}
	if input.Email == "" {
		errorsSlice = append(errorsSlice, "email is required")
	}
	if input.Phone == "" {
		errorsSlice = append(errorsSlice, "phone is required")
	}
	if input.RequestedAmount <= 0 {
		errorsSlice = append(errorsSlice, "requested amount must be greater than zero")
	}
	if input.TermMonths <= 0 {
		errorsSlice = append(errorsSlice, "term months must be greater than zero")
	}
	if len(errorsSlice) > 0 {
		return errors.New("validation errors: " + strings.Join(errorsSlice, ", "))
	}
	return nil
}
