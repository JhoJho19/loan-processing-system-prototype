package service

import (
	"loan-processing-system/internal/application/domain"
	"testing"
)

func TestApplicationService_CreateApplication(t *testing.T) {
	//Arrenge

	//Act

	//Assert
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
