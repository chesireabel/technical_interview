package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	// Loading .env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// new connection pool
	DB, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	// Test the connection
	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	fmt.Println("Database connection successful!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed")
	}
}