package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		fmt.Println("DATABASE_URL is not configured")
		return
	}

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
