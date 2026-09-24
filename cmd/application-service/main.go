package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := "postgres://loan_user:loan_password@localhost:5432/loan_processing"

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Println("Error creating connection pool:", err)
		return
	}

	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}

	fmt.Println("Successfully connected to PostgreSQL")
}
