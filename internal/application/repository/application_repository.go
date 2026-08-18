package repository

import "loan-processing-system/internal/application/domain"

type ApplicationRepository interface {
	Create(customer domain.Customer, application domain.LoanApplication) error
	GetByID(id string) (domain.LoanApplication, error)
}
