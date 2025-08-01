package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Linux; Android 14; SM-S918B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.7204.49 Safari/537.36",
	"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.7103.125 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.7204.96 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.7204.96 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_7_5 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/138.0.7204.119 Mobile/15E148 Safari/604.1",
}

func getRandomUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	return userAgents[rand.Intn(len(userAgents))]
}

type Chat struct {
	ChatId       int
	Time         string
	TabulationId int
}

func PegarChats(cid int, apikey string) ([]Chat, error) {
	log.Printf("Starting to fetch chats for user ID: %d", cid)

	var stop_flag = false
	contador := 0
	auth := fmt.Sprintf("Bearer %s", apikey)
	var filteredChats []Chat

	log.Println("Beginning pagination loop to fetch all chats...")

	for stop_flag != true {
		url := fmt.Sprintf("https://api.huggy.app/v3/chats?agent=%d&situation=in_chat&page=%d", cid, contador)
		log.Printf("Making API request to page %d: %s", contador, url)

		contador++
		httpReq, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("ERROR: Failed to create HTTP request: %v", err)
			return nil, err
		}

		userAgent := getRandomUserAgent()
		log.Printf("Using User-Agent: %s", userAgent)

		httpReq.Header.Set("Authorization", auth)
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("User-Agent", userAgent)

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		if err != nil {
			log.Printf("ERROR: HTTP request failed: %v", err)
			return nil, err
		}
		defer resp.Body.Close()

		log.Printf("Received response with status: %s", resp.Status)

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			log.Printf("ERROR: Request failed with status: %s", resp.Status)
			return nil, fmt.Errorf("request failed with status: %s", resp.Status)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("ERROR: Failed to read response body: %v", err)
			return nil, err
		}

		log.Printf("Response body size: %d bytes", len(bodyBytes))

		type ChatResp struct {
			ID             int `json:"id"`
			ChatTabulation struct {
				ID string `json:"id"`
			} `json:"chatTabulation"`
			LastMessage struct {
				SendAt string `json:"sendAt"`
			} `json:"lastMessage"`
		}

		var chats []ChatResp
		err = json.Unmarshal(bodyBytes, &chats)
		if err != nil {
			log.Printf("ERROR: Failed to unmarshal response: %v", err)
			log.Printf("Response body: %s", string(bodyBytes))
		} else {
			log.Printf("SUCCESS: Received %d chats on page %d", len(chats), contador-1)
			if len(chats) == 0 {
				log.Println("No more chats found, stopping pagination")
				stop_flag = true
			}

			validChats := 0
			for _, chat := range chats {
				if chat.LastMessage.SendAt != "" {
					t, err := time.Parse("2006-01-02 15:04:05", chat.LastMessage.SendAt)
					cutoff := time.Now().AddDate(0, 0, -3)
					is_available := t.Before(cutoff) || t.Equal(cutoff)
					if err == nil && is_available {
						tabulationId := parseTabId(chat.ChatTabulation.ID)
						log.Printf("Processing chat %d - TabulationId: %d (original: '%s')", chat.ID, tabulationId, chat.ChatTabulation.ID)

						filteredChats = append(filteredChats, Chat{
							ChatId:       chat.ID,
							Time:         chat.LastMessage.SendAt,
							TabulationId: tabulationId,
						})
						validChats++
					} else {
						log.Printf("Skipping chat %d - Invalid time format or after cutoff. Time: %s", chat.ID, chat.LastMessage.SendAt)
					}
				} else {
					log.Printf("Skipping chat %d - No last message time", chat.ID)
				}
			}
			log.Printf("Added %d valid chats from page %d", validChats, contador-1)
		}

		log.Printf("Waiting 1 second before next request...")
		time.Sleep(time.Second)
	}

	log.Printf("SUCCESS: Completed fetching chats for user %d", cid)
	log.Printf("Total filtered chats: %d", len(filteredChats))

	for _, chat := range filteredChats {
		log.Printf("Filtered chat - ID: %d, TabulationId: %d, Time: %s", chat.ChatId, chat.TabulationId, chat.Time)
	}

	return filteredChats, nil
}

func parseTabId(id string) int {
	if strings.TrimSpace(id) == "" {
		log.Printf("INFO: Empty tabulation ID found, using 0 as default")
		return 0
	}

	tabId, err := strconv.Atoi(strings.TrimSpace(id))
	if err != nil {
		log.Printf("WARNING: Failed to parse tabulation ID '%s': %v - using 0 as default", id, err)
		return 0
	}

	return tabId
}

func FecharChat(chat Chat, apikey string) error {
	log.Printf("Starting to close chat ID: %d", chat.ChatId)

	url := fmt.Sprintf("https://api.huggy.app/v3/chats/%d/close", chat.ChatId)
	log.Printf("API endpoint: %s", url)

	payload := map[string]interface{}{
		"tabulation":   chat.TabulationId,
		"comment":      "Fechado via robô.",
		"sendFeedback": false,
	}

	log.Printf("Request payload: tabulation=%d, comment='Fechado via robô.', sendFeedback=false", chat.TabulationId)

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Printf("ERROR: Failed to marshal JSON payload: %v", err)
		return err
	}

	httpReq, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("ERROR: Failed to create HTTP request: %v", err)
		return err
	}

	userAgent := getRandomUserAgent()
	log.Printf("Using User-Agent: %s", userAgent)

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))
	httpReq.Header.Set("User-Agent", userAgent)

	client := &http.Client{}
	resp, err := client.Do(httpReq)

	if err != nil {
		log.Printf("ERROR: HTTP request failed: %v", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("Received response with status: %s", resp.Status)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		log.Printf("ERROR: Request failed with status: %s", resp.Status)
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	log.Printf("SUCCESS: Chat %d closed successfully", chat.ChatId)
	return nil
}
