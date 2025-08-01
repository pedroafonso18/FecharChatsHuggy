package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func FetchUsers(db *sql.DB) ([]string, error) {
	log.Println("Executing query to fetch users with cargo='vendedor'...")

	query := `SELECT userid FROM usuarios WHERE cargo = 'vendedor' `
	log.Printf("SQL Query: %s", query)

	var vendedores []string
	res, err := db.Query(query)
	if err != nil {
		log.Printf("ERROR: Database query failed: %v", err)
		return nil, err
	}
	defer res.Close()

	log.Println("Query executed successfully, scanning results...")

	rowCount := 0
	for res.Next() {
		var vendedor string
		err := res.Scan(&vendedor)
		if err != nil {
			log.Printf("ERROR: Failed to scan row %d: %v", rowCount+1, err)
			return nil, fmt.Errorf("failed to scan instance row: %w", err)
		}
		vendedores = append(vendedores, vendedor)
		rowCount++

		if rowCount%10 == 0 {
			log.Printf("Processed %d users so far...", rowCount)
		}
	}

	if err = res.Err(); err != nil {
		log.Printf("ERROR: Error during result iteration: %v", err)
		return nil, fmt.Errorf("error iterating instances: %w", err)
	}

	log.Printf("SUCCESS: Retrieved %d users from database", len(vendedores))
	return vendedores, nil
}
