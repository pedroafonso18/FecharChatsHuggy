package main

import (
	internal "FecharChats/internal/api"
	"FecharChats/internal/config"
	"FecharChats/internal/database"
	"log"
	"strconv"
	"time"
)

func main() {
	log.Println("=== Starting FecharChats Application ===")

	cfg := config.LoadEnv()
	log.Printf("Configuration loaded - DB URL: %s, API Key: %s",
		maskString(cfg.DbUrl, 10), maskString(cfg.ApiKey, 8))

	for {
		log.Println("--- Starting new iteration ---")

		log.Println("Attempting to connect to database...")
		db, err := database.ConnectDb(cfg.DbUrl)
		if err != nil {
			log.Printf("ERROR: Failed to connect to database: %v", err)
			log.Println("Waiting 5 minutes before retry...")
			time.Sleep(5 * time.Minute)
			continue
		}
		log.Println("SUCCESS: Database connection established")

		log.Println("Fetching users from database...")
		fetch, err := database.FetchUsers(db)
		if err != nil {
			log.Printf("ERROR: Failed to fetch users: %v", err)
			log.Println("Waiting 5 minutes before retry...")
			time.Sleep(5 * time.Minute)
			continue
		}
		log.Printf("SUCCESS: Retrieved %d users from database", len(fetch))

		processedChats := 0
		for i, userStr := range fetch {
			log.Printf("Processing user %d/%d: %s", i+1, len(fetch), userStr)

			userID, err := strconv.Atoi(userStr)
			if err != nil {
				log.Printf("ERROR: Failed to convert user ID '%s' to int: %v", userStr, err)
				continue
			}

			chats, err := internal.PegarChats(userID, cfg.ApiKey)
			if err != nil {
				log.Printf("ERROR: Failed to fetch chats for user %s: %v", userStr, err)
				continue
			}

			log.Printf("Found %d chats for user %s", len(chats), userStr)

			for j, chat := range chats {
				log.Printf("Processing chat %d/%d for user %s - ID: %d, TabulationId: %d, Time: %s",
					j+1, len(chats), userStr, chat.ChatId, chat.TabulationId, chat.Time)

				chatLastMessage, err := time.Parse("2006-01-02 15:04:05", chat.Time)
				if err != nil {
					log.Printf("ERROR: Failed to parse chat time '%s' for chat %d: %v", chat.Time, chat.ChatId, err)
					chatLastMessage = time.Time{}
				}

				log.Printf("Logging chat %d to database before API call", chat.ChatId)
				err = database.InsertLog(chat.ChatId, userID, chat.TabulationId, chatLastMessage, db)
				if err != nil {
					log.Printf("ERROR: Failed to log chat %d to database: %v", chat.ChatId, err)
					continue
				}
				log.Printf("SUCCESS: Chat %d logged to database", chat.ChatId)

				err = internal.FecharChat(chat, cfg.ApiKey)
				if err != nil {
					log.Printf("ERROR: Failed to close chat %d via API: %v", chat.ChatId, err)
					log.Printf("NOTE: Chat %d was already logged to database, so we have a record", chat.ChatId)
				} else {
					log.Printf("SUCCESS: Chat %d closed successfully via API", chat.ChatId)
				}

				processedChats++

				time.Sleep(time.Second)
			}
		}

		log.Printf("=== Iteration complete - Processed %d chats ===", processedChats)
		db.Close()
		log.Println("Waiting before next iteration...")
		time.Sleep(30 * time.Second)
	}
}

func maskString(s string, visibleChars int) string {
	if len(s) <= visibleChars {
		return "***"
	}
	return s[:visibleChars] + "***"
}
