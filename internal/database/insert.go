package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

func InsertLog(cid, uid, tid int, chatLastMessage time.Time, db *sql.DB) error {
	log.Printf("Starting log insertion - ChatId: %d, UserId: %d, TabulationId: %d, LastMessage: %s", cid, uid, tid, chatLastMessage.Format("2006-01-02 15:04:05"))

	if cid <= 0 {
		log.Printf("ERROR: Invalid ChatId %d - must be greater than 0", cid)
		return fmt.Errorf("invalid ChatId: %d", cid)
	}

	if uid <= 0 {
		log.Printf("ERROR: Invalid UserId %d - must be greater than 0", uid)
		return fmt.Errorf("invalid UserId: %d", uid)
	}

	if tid <= 0 {
		log.Printf("WARNING: TabulationId is %d (null/zero) - this may indicate an issue with the chat data", tid)
		log.Printf("Proceeding with insertion but tabulation_id will be %d", tid)
	}

	// Validate timestamp
	if chatLastMessage.IsZero() {
		log.Printf("WARNING: Chat last message timestamp is zero - this may indicate missing time data")
	}

	log.Printf("Validation passed - proceeding with database insertion")

	query := `
        INSERT INTO logs_fechamento_chamados (chatid, userid, tabulation_id, chat_last_message)
        VALUES ($1, $2, $3, $4)
    `
	log.Printf("Executing SQL query: %s", query)
	log.Printf("Query parameters: ChatId=%d, UserId=%d, TabulationId=%d, LastMessage=%s", cid, uid, tid, chatLastMessage.Format("2006-01-02 15:04:05"))

	result, err := db.Exec(query, cid, uid, tid, chatLastMessage)
	if err != nil {
		log.Printf("ERROR: Database insertion failed for ChatId %d: %v", cid, err)
		return fmt.Errorf("insert failed: %w", err)
	}

	// Get the number of affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("WARNING: Could not determine rows affected: %v", err)
	} else {
		log.Printf("SUCCESS: Inserted %d row(s) for ChatId %d", rowsAffected, cid)
	}

	log.Printf("SUCCESS: Log insertion completed for ChatId: %d", cid)
	return nil
}
