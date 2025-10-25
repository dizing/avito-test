package main

import (
	"avito-test/internal/handler"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

const BASE_URL = "http://localhost:8080"

func register(userNumber int) (string, error) {
	username := fmt.Sprintf("vu_%d_user", userNumber)
	authPayload := handler.AuthRequest{
		Username: username,
		Password: "pass",
	}

	jsonData, err := json.Marshal(authPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth payload: %v", err)
	}

	resp, err := http.Post(BASE_URL+"/api/auth", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp handler.AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %v", err)
	}

	return authResp.Token, nil
}

func worker(workerID int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for userNumber := range jobs {
		token, err := register(userNumber)
		if err != nil {
			log.Printf("ERROR: Worker %d failed to register user %d: %v", workerID, userNumber, err)
			continue
		}

		if userNumber%100 == 0 {
			log.Printf("SUCCESS: Worker %d registered user %d, token: %s", workerID, userNumber, token)
		}
	}
}

func main() {
	// Configuration
	numUsers := 300 * 100
	numWorkers := 50

	// Create channels and wait group
	jobs := make(chan int, numUsers)
	var wg sync.WaitGroup

	// Start worker pool
	log.Printf("Starting %d workers to register %d users", numWorkers, numUsers)

	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, jobs, &wg)
	}

	// Send jobs to workers
	startTime := time.Now()

	for i := 1; i <= numUsers; i++ {
		jobs <- i
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()

	duration := time.Since(startTime)
	log.Printf("Completed all registrations in %v", duration)
}
