package service

import (
	"loan-processing-system/internal/application/domain"
	"time"
)

type ApplicationService struct{}

func (c ApplicationService) CreateApplication(customerApplicationInput CreateApplicationInput) (domain.LoanApplication, error) {

	loanApplication := domain.LoanApplication{}

	// нужно добавить будет логику генерации айдишников
	// нужно будет перебирать айдишники и проверять чтобы не было повтора

	loanApplication.RequestedAmount = customerApplicationInput.RequestedAmount
	loanApplication.TermMonths = customerApplicationInput.TermMonths
	loanApplication.Status = "New"
	loanApplication.CreatedAt = time.Now()

	return loanApplication, nil
}

func (c ApplicationService) GetApplication(applicationID int) (domain.LoanApplication, error) {
	// тут нужно будет выгрузить клиентов и перебрать по ID
	return domain.LoanApplication{}, nil
}
