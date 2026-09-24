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

	defer tx.Rollback((ctx))

	return nil
}

func (r *PostgresRepository) GetByID(
	applicationID string,
) (domain.LoanApplication, error) {
	return domain.LoanApplication{}, nil
}
