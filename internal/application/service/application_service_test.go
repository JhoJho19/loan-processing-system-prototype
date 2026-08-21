package service

import (
	"loan-processing-system/internal/application/domain"
	"testing"
	"time"
)

func TestApplicationService_CreateApplication(t *testing.T) {
	// 1. Arrange (Подготовка)
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
			t.Errorf("Expected FirstName %s, got %s", inputFirstName, customer.FirstName)
		}
		if customer.LastName != inputLastName {
			t.Errorf("Expected LastName %s, got %s", inputLastName, customer.LastName)
		}
		if customer.Address != inputAddress {
			t.Errorf("Expected Address %s, got %s", inputAddress, customer.Address)
		}
		if customer.BirthDate != inputBirthDate {
			t.Errorf("Expected BirthDate %s, got %s", inputBirthDate, customer.BirthDate)
		}
		if customer.MonthlyIncome != inputMonthlyIncome {
			t.Errorf("Expected MonthlyIncome %v, got %d", inputMonthlyIncome, customer.MonthlyIncome)
		}
		if customer.EmploymentType != inputEmploymentType {
			t.Errorf("Expected EmploymentType %s, got %s", inputEmploymentType, customer.EmploymentType)
		}
		if customer.Email != inputEmail {
			t.Errorf("Expected Email %s, got %s", inputEmail, customer.Email)
		}
		if customer.Phone != inputPhone {
			t.Errorf("Expected Phone %s, got %s", inputPhone, customer.Phone)
		}
		if application.RequestedAmount != inputRequestedAmount {
			t.Errorf("Expected RequestedAmount %v, got %d", inputRequestedAmount, application.RequestedAmount)
		}
		if application.TermMonths != inputTermMonths {
			t.Errorf("Expected TermMonths %d, got %d", inputTermMonths, application.TermMonths)
		}
		return nil
	}

	service := NewApplicationService(fakeRepo)

	// 2. Act (Действие)
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

	createdApplication, err := service.CreateApplication(applicationInput)

	// 3. Assert (Проверка)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if createdApplication.ID == "" {
		t.Errorf("Expected created application to have an ID")
	}
	if createdApplication.CustomerID == "" {
		t.Errorf("Expected created application to have a CustomerID")
	}
	if createdApplication.RequestedAmount != int(inputRequestedAmount) {
		t.Errorf("Expected RequestedAmount %v, got %d", inputRequestedAmount, createdApplication.RequestedAmount)
	}
	if createdApplication.TermMonths != inputTermMonths {
		t.Errorf("Expected TermMonths %d, got %d", inputTermMonths, createdApplication.TermMonths)
	}
}

func TestApplicationService_GetApplication(t *testing.T) {
	//Arrenge

	//Act

	//Assert
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
