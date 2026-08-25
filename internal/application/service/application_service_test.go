package service

import (
	"errors"
	"loan-processing-system/internal/application/domain"
	"strings"
	"testing"
	"time"
)

func TestApplicationService_CreateApplication(t *testing.T) {
	// 1. Arrange
	var (
		inputFirstName       = "John"
		inputLastName        = "Doe"
		inputAddress         = "123 Main St"
		inputBirthDate       = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		inputMonthlyIncome   = 5000
		inputEmploymentType  = "Full-time"
		inputEmail           = "john.doe@example.com"
		inputPhone           = "123-456-7890"
		inputRequestedAmount = 10000
		inputTermMonths      = 12
	)

	fakeRepo := &fakeApplicationRepository{}

	fakeRepo.CreateFunc = func(customer domain.Customer, application domain.LoanApplication) error {
		if customer.FirstName != inputFirstName {
			t.Errorf("expected FirstName %s, got %s", inputFirstName, customer.FirstName)
		}
		if customer.LastName != inputLastName {
			t.Errorf("expected LastName %s, got %s", inputLastName, customer.LastName)
		}
		if customer.Address != inputAddress {
			t.Errorf("expected Address %s, got %s", inputAddress, customer.Address)
		}
		if customer.BirthDate != inputBirthDate {
			t.Errorf("expected BirthDate %v, got %v", inputBirthDate, customer.BirthDate)
		}
		if customer.MonthlyIncome != inputMonthlyIncome {
			t.Errorf("expected MonthlyIncome %d, got %d", inputMonthlyIncome, customer.MonthlyIncome)
		}
		if customer.EmploymentType != inputEmploymentType {
			t.Errorf("expected EmploymentType %s, got %s", inputEmploymentType, customer.EmploymentType)
		}
		if customer.Email != inputEmail {
			t.Errorf("expected Email %s, got %s", inputEmail, customer.Email)
		}
		if customer.Phone != inputPhone {
			t.Errorf("expected Phone %s, got %s", inputPhone, customer.Phone)
		}
		if application.RequestedAmount != inputRequestedAmount {
			t.Errorf("expected RequestedAmount %d, got %d", inputRequestedAmount, application.RequestedAmount)
		}
		if application.TermMonths != inputTermMonths {
			t.Errorf("expected TermMonths %d, got %d", inputTermMonths, application.TermMonths)
		}

		return nil
	}

	service := NewApplicationService(fakeRepo)

	applicationInput := CreateApplicationInput{
		FirstName:       inputFirstName,
		LastName:        inputLastName,
		Address:         inputAddress,
		BirthDate:       inputBirthDate,
		MonthlyIncome:   inputMonthlyIncome,
		EmploymentType:  inputEmploymentType,
		Email:           inputEmail,
		Phone:           inputPhone,
		RequestedAmount: inputRequestedAmount,
		TermMonths:      inputTermMonths,
	}

	// 2. Act
	createdApplication, err := service.CreateApplication(applicationInput)

	// 3. Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdApplication.ID == "" {
		t.Errorf("expected created application to have an ID")
	}

	if createdApplication.CustomerID == "" {
		t.Errorf("expected created application to have a CustomerID")
	}

	if createdApplication.RequestedAmount != inputRequestedAmount {
		t.Errorf(
			"expected RequestedAmount %d, got %d",
			inputRequestedAmount,
			createdApplication.RequestedAmount,
		)
	}

	if createdApplication.TermMonths != inputTermMonths {
		t.Errorf(
			"expected TermMonths %d, got %d",
			inputTermMonths,
			createdApplication.TermMonths,
		)
	}
}

func TestApplicationService_CreateApplication_InvalidInput(t *testing.T) {
	// 1. Arrange
	var (
		inputFirstName       = ""
		inputLastName        = "Doe"
		inputAddress         = "123 Main St"
		inputBirthDate       = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		inputMonthlyIncome   = 5000
		inputEmploymentType  = "Full-time"
		inputEmail           = "john.doe@example.com"
		inputPhone           = "123-456-7890"
		inputRequestedAmount = 10000
		inputTermMonths      = 12
	)

	applicationInput := CreateApplicationInput{
		FirstName:       inputFirstName,
		LastName:        inputLastName,
		Address:         inputAddress,
		BirthDate:       inputBirthDate,
		MonthlyIncome:   inputMonthlyIncome,
		EmploymentType:  inputEmploymentType,
		Email:           inputEmail,
		Phone:           inputPhone,
		RequestedAmount: inputRequestedAmount,
		TermMonths:      inputTermMonths,
	}

	createCalled := false

	fakeRepo := &fakeApplicationRepository{
		CreateFunc: func(_ domain.Customer, _ domain.LoanApplication) error {
			createCalled = true
			return nil
		},
	}

	service := NewApplicationService(fakeRepo)

	// 2. Act
	createdApplication, err := service.CreateApplication(applicationInput)

	// 3. Assert
	if err == nil {
		t.Fatalf("expected error for invalid input, got nil")
	}

	if !strings.Contains(err.Error(), "first name is required") {
		t.Errorf("expected error message about first name, got %v", err)
	}

	if createCalled {
		t.Errorf("expected repository Create not to be called for invalid input")
	}

	if createdApplication != (domain.LoanApplication{}) {
		t.Errorf("expected empty LoanApplication, got %+v", createdApplication)
	}
}

func TestApplicationService_CreateApplication_RepositoryError(t *testing.T) {
	// 1. Arrange
	var (
		inputFirstName       = "John"
		inputLastName        = "Doe"
		inputAddress         = "123 Main St"
		inputBirthDate       = time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		inputMonthlyIncome   = 5000
		inputEmploymentType  = "Full-time"
		inputEmail           = "john.doe@example.com"
		inputPhone           = "123-456-7890"
		inputRequestedAmount = 10000
		inputTermMonths      = 12
	)

	fakeRepo := &fakeApplicationRepository{}

	repositoryError := errors.New("repository error")

	createCalled := false

	fakeRepo.CreateFunc = func(customer domain.Customer, application domain.LoanApplication) error {
		createCalled = true
		return repositoryError
	}

	service := NewApplicationService(fakeRepo)

	applicationInput := CreateApplicationInput{
		FirstName:       inputFirstName,
		LastName:        inputLastName,
		Address:         inputAddress,
		BirthDate:       inputBirthDate,
		MonthlyIncome:   inputMonthlyIncome,
		EmploymentType:  inputEmploymentType,
		Email:           inputEmail,
		Phone:           inputPhone,
		RequestedAmount: inputRequestedAmount,
		TermMonths:      inputTermMonths,
	}

	// 2. Act

	createdApplication, err := service.CreateApplication(applicationInput)

	// 3. Assert
	if createCalled == false {
		t.Errorf("expected repository Create to be called")
	}
	if err == nil {
		t.Fatalf("expected error from repository, got nil")
	}
	if err != repositoryError {
		t.Errorf("expected error %v, got %v", repositoryError, err)
	}
	if createdApplication != (domain.LoanApplication{}) {
		t.Errorf("expected empty LoanApplication, got %+v", createdApplication)
	}
}

func TestApplicationService_GetApplication(t *testing.T) {
	// 1. Arrange
	applicationID := "some-application-id"

	expectedApplication := domain.LoanApplication{
		ID:              applicationID,
		CustomerID:      "some-customer-id",
		RequestedAmount: 10000,
		TermMonths:      12,
		Status:          domain.StatusNew,
	}

	getByIDCalled := false

	fakeRepo := &fakeApplicationRepository{
		GetByIDFunc: func(id string) (domain.LoanApplication, error) {
			getByIDCalled = true

			if id != applicationID {
				t.Errorf("expected GetByID to be called with %s, got %s", applicationID, id)
			}

			return expectedApplication, nil
		},
	}

	service := NewApplicationService(fakeRepo)

	// 2. Act
	application, err := service.GetApplication(applicationID)

	// 3. Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !getByIDCalled {
		t.Errorf("expected repository GetByID to be called")
	}

	if application != expectedApplication {
		t.Errorf("expected application %+v, got %+v", expectedApplication, application)
	}
}

func TestApplicationService_GetApplication_InvalidID(t *testing.T) {
	// 1. Arrange
	applicationID := ""

	getByIDCalled := false

	fakeRepo := &fakeApplicationRepository{
		GetByIDFunc: func(id string) (domain.LoanApplication, error) {
			getByIDCalled = true
			return domain.LoanApplication{}, nil
		},
	}

	service := NewApplicationService(fakeRepo)

	// 2. Act
	application, err := service.GetApplication(applicationID)

	// 3. Assert
	if err == nil {
		t.Fatalf("expected error for invalid application ID, got nil")
	}

	if err != domain.ErrInvalidApplicationID {
		t.Errorf(
			"expected error %v, got %v",
			domain.ErrInvalidApplicationID,
			err,
		)
	}

	if getByIDCalled {
		t.Errorf("expected repository GetByID not to be called")
	}

	if application != (domain.LoanApplication{}) {
		t.Errorf("expected empty LoanApplication, got %+v", application)
	}
}

func TestApplicationService_GetApplication_RepositoryError(t *testing.T) {
	// 1. Arrange
	applicationID := "some-application-id"
	repositoryError := errors.New("repository error")
	getByIDCalled := false

	fakeRepo := &fakeApplicationRepository{
		GetByIDFunc: func(id string) (domain.LoanApplication, error) {
			getByIDCalled = true
			return domain.LoanApplication{}, repositoryError
		},
	}

	service := NewApplicationService(fakeRepo)

	// 2. Act
	application, err := service.GetApplication(applicationID)

	// 3. Assert
	if err == nil {
		t.Fatalf("expected error from repository, got nil")
	}
	if err != repositoryError {
		t.Errorf("expected error %v, got %v", repositoryError, err)
	}
	if !getByIDCalled {
		t.Errorf("expected repository GetByID to be called")
	}
	if application != (domain.LoanApplication{}) {
		t.Errorf("expected empty LoanApplication, got %+v", application)
	}
}

type fakeApplicationRepository struct {
	CreateFunc  func(customer domain.Customer, application domain.LoanApplication) error
	GetByIDFunc func(id string) (domain.LoanApplication, error)
}

func (f *fakeApplicationRepository) Create(customer domain.Customer, application domain.LoanApplication) error {
	if f.CreateFunc != nil {
		return f.CreateFunc(customer, application)
	}

	return nil
}

func (f *fakeApplicationRepository) GetByID(id string) (domain.LoanApplication, error) {
	if f.GetByIDFunc != nil {
		return f.GetByIDFunc(id)
	}

	return domain.LoanApplication{}, nil
}
