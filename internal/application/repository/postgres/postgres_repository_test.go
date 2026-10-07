package postgres

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"loan-processing-system/internal/application/domain"
)

func setupPostgres(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(20*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres: %v", err)
	}
	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Logf("failed to stop postgres: %v", err)
		}
	})

	connString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		t.Fatalf("failed to create pgx pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("postgres is unavailable: %v", err)
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate test source file")
	}
	migrationPath := filepath.Join(filepath.Dir(filename), "../../../../migrations/000001_init.up.sql")
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		t.Fatalf("failed to read database migration: %v", err)
	}
	if _, err := pool.Exec(ctx, string(migration)); err != nil {
		t.Fatalf("failed to apply database migration: %v", err)
	}

	return pool
}

func testCustomer(id string) domain.Customer {
	return domain.Customer{
		ID:             id,
		FirstName:      "John",
		LastName:       "Doe",
		BirthDate:      time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		MonthlyIncome:  5000,
		EmploymentType: string(domain.StatusEmployed),
		Email:          "john.doe@example.com",
		Phone:          "+1234567890",
		Address:        "123 Main St",
	}
}

func testApplication(id, customerID string) domain.LoanApplication {
	now := time.Now().UTC().Truncate(time.Microsecond)
	return domain.LoanApplication{
		ID:              id,
		CustomerID:      customerID,
		RequestedAmount: 10000,
		TermMonths:      12,
		Status:          domain.StatusNew,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

func TestPostgresRepository_Create_Success(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	customer := testCustomer("00000000-0000-0000-0000-000000000001")
	application := testApplication("00000000-0000-0000-0000-000000000101", customer.ID)
	ctx := context.Background()

	if err := repo.Create(customer, application); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	var gotCustomer domain.Customer
	if err := pool.QueryRow(ctx, `
		SELECT id, first_name, last_name, birth_date, monthly_income,
		       employment_type, email, phone, address
		FROM customers WHERE id = $1`, customer.ID,
	).Scan(
		&gotCustomer.ID,
		&gotCustomer.FirstName,
		&gotCustomer.LastName,
		&gotCustomer.BirthDate,
		&gotCustomer.MonthlyIncome,
		&gotCustomer.EmploymentType,
		&gotCustomer.Email,
		&gotCustomer.Phone,
		&gotCustomer.Address,
	); err != nil {
		t.Fatalf("failed to read inserted customer: %v", err)
	}
	if gotCustomer.ID != customer.ID ||
		gotCustomer.FirstName != customer.FirstName ||
		gotCustomer.LastName != customer.LastName ||
		!gotCustomer.BirthDate.Equal(customer.BirthDate) ||
		gotCustomer.MonthlyIncome != customer.MonthlyIncome ||
		gotCustomer.EmploymentType != customer.EmploymentType ||
		gotCustomer.Email != customer.Email ||
		gotCustomer.Phone != customer.Phone ||
		gotCustomer.Address != customer.Address {
		t.Errorf("stored customer = %+v, want %+v", gotCustomer, customer)
	}

	var gotApplication domain.LoanApplication
	if err := pool.QueryRow(ctx, `
		SELECT id, customer_id, requested_amount, term_months, status, created_at, updated_at
		FROM loan_applications WHERE id = $1`, application.ID,
	).Scan(
		&gotApplication.ID,
		&gotApplication.CustomerID,
		&gotApplication.RequestedAmount,
		&gotApplication.TermMonths,
		&gotApplication.Status,
		&gotApplication.CreatedAt,
		&gotApplication.UpdatedAt,
	); err != nil {
		t.Fatalf("failed to read inserted application: %v", err)
	}
	assertApplicationEqual(t, gotApplication, application)
}

func TestPostgresRepository_Create_RollbackIfApplicationInsertFails(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	customer := testCustomer("00000000-0000-0000-0000-000000000002")
	application := testApplication("00000000-0000-0000-0000-000000000102", customer.ID)
	application.TermMonths = -1

	if err := repo.Create(customer, application); err == nil {
		t.Fatal("Create() error = nil, want term_months constraint violation")
	}

	assertRowCount(t, pool, "customers", customer.ID, 0)
	assertRowCount(t, pool, "loan_applications", application.ID, 0)
}

func TestPostgresRepository_Create_DuplicateCustomerID_ReturnsError(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	customer := testCustomer("00000000-0000-0000-0000-000000000003")
	firstApplication := testApplication("00000000-0000-0000-0000-000000000103", customer.ID)
	secondApplication := testApplication("00000000-0000-0000-0000-000000000203", customer.ID)

	if err := repo.Create(customer, firstApplication); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	if err := repo.Create(customer, secondApplication); err == nil {
		t.Fatal("second Create() error = nil, want primary key violation")
	}

	assertRowCount(t, pool, "customers", customer.ID, 1)
	assertRowCount(t, pool, "loan_applications", secondApplication.ID, 0)
}

func TestPostgresRepository_GetByID_Success(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	customer := testCustomer("00000000-0000-0000-0000-000000000004")
	application := testApplication("00000000-0000-0000-0000-000000000104", customer.ID)

	if err := repo.Create(customer, application); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	got, err := repo.GetByID(application.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	assertApplicationEqual(t, got, application)
}

func TestPostgresRepository_GetByID_NotFound(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)

	got, err := repo.GetByID("00000000-0000-0000-0000-000000000999")
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("GetByID() error = %v, want %v", err, pgx.ErrNoRows)
	}
	if got != (domain.LoanApplication{}) {
		t.Errorf("GetByID() result = %+v, want zero value", got)
	}
}

func TestPostgresRepository_GetByID_ReturnsCorrectApplication(t *testing.T) {
	pool := setupPostgres(t)
	repo := NewPostgresRepository(pool)
	applications := []domain.LoanApplication{
		testApplication("00000000-0000-0000-0000-000000000105", "00000000-0000-0000-0000-000000000005"),
		testApplication("00000000-0000-0000-0000-000000000205", "00000000-0000-0000-0000-000000000006"),
		testApplication("00000000-0000-0000-0000-000000000305", "00000000-0000-0000-0000-000000000007"),
	}
	for _, application := range applications {
		if err := repo.Create(testCustomer(application.CustomerID), application); err != nil {
			t.Fatalf("Create(%q) error = %v", application.ID, err)
		}
	}

	got, err := repo.GetByID(applications[1].ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	assertApplicationEqual(t, got, applications[1])
}

func assertRowCount(t *testing.T, pool *pgxpool.Pool, table, id string, want int) {
	t.Helper()

	query := "SELECT COUNT(*) FROM " + table + " WHERE id = $1"
	var got int
	if err := pool.QueryRow(context.Background(), query, id).Scan(&got); err != nil {
		t.Fatalf("failed to count %s rows: %v", table, err)
	}
	if got != want {
		t.Errorf("%s row count for %q = %d, want %d", table, id, got, want)
	}
}

func assertApplicationEqual(t *testing.T, got, want domain.LoanApplication) {
	t.Helper()

	if got.ID != want.ID {
		t.Errorf("ID = %q, want %q", got.ID, want.ID)
	}
	if got.CustomerID != want.CustomerID {
		t.Errorf("CustomerID = %q, want %q", got.CustomerID, want.CustomerID)
	}
	if got.RequestedAmount != want.RequestedAmount {
		t.Errorf("RequestedAmount = %d, want %d", got.RequestedAmount, want.RequestedAmount)
	}
	if got.TermMonths != want.TermMonths {
		t.Errorf("TermMonths = %d, want %d", got.TermMonths, want.TermMonths)
	}
	if got.Status != want.Status {
		t.Errorf("Status = %q, want %q", got.Status, want.Status)
	}
	assertTimeEqual(t, "CreatedAt", got.CreatedAt, want.CreatedAt)
	assertTimeEqual(t, "UpdatedAt", got.UpdatedAt, want.UpdatedAt)
}

func assertTimeEqual(t *testing.T, name string, got, want time.Time) {
	t.Helper()

	got = got.UTC().Truncate(time.Microsecond)
	want = want.UTC().Truncate(time.Microsecond)
	if !got.Equal(want) {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}
