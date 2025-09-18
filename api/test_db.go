package main

import (
	"fmt"
	"log"
	"github.com/biplab-sutradhar/slugify/api/internal/db"
)

func main() {
	connStr := "postgres://postgres:mynew@localhost:5432/urlshortener?sslmode=disable"
	database, err := db.NewDB(connStr)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	var version string
	err = database.QueryRow("SELECT version()").Scan(&version)
	if err != nil {
		log.Fatalf("Failed to query version: %v", err)
	}
	fmt.Println("Connected! PostgreSQL version:", version)
}
