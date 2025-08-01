package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func ConnectDb(dbURL string) (*sql.DB, error) {
	log.Printf("Attempting to connect to database...")
	log.Printf("Database URL: %s", maskDbURL(dbURL))

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("ERROR: Failed to open database connection: %v", err)
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	log.Println("Database connection opened, testing connection...")

	if err := db.Ping(); err != nil {
		log.Printf("ERROR: Database ping failed: %v", err)
		db.Close()
		return nil, fmt.Errorf("connection to DB failed: %w", err)
	}

	log.Println("SUCCESS: Database connection established and ping successful")
	return db, nil
}

func maskDbURL(url string) string {
	if len(url) <= 20 {
		return "***"
	}
	return url[:10] + "***" + url[len(url)-10:]
}
