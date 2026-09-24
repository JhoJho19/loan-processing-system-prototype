package postgres

import (
	"context"
	"loan-processing-system/internal/application/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

func (r *PostgresRepository) Create(
	customer domain.Customer,
	application domain.LoanApplication,
) error {

	ctx := context.Background()

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`INSERT INTO customers (
			id,
			first_name,
			last_name,
			birth_date,
			monthly_income,
			employment_type,
			email,
			phone,
			address
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		customer.ID,
		customer.FirstName,
		customer.LastName,
		customer.BirthDate,
		customer.MonthlyIncome,
		customer.EmploymentType,
		customer.Email,
		customer.Phone,
		customer.Address,
	)

	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO loan_applications (
			id,
			customer_id,
			requested_amount,
			term_months,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		application.ID,
		application.CustomerID,
		application.RequestedAmount,
		application.TermMonths,
		application.Status,
		application.CreatedAt,
		application.UpdatedAt,
	)

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByID(
	applicationID string,
) (domain.LoanApplication, error) {
	return domain.LoanApplication{}, nil
}
