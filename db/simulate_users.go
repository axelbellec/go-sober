package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"go-sober/platform"
)

func init() {
	platform.InitPlatform()
}

type AuthResponse struct {
	Token string `json:"token"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DrinkLog struct {
	DrinkOptionID int       `json:"drink_option_id"`
	LoggedAt      time.Time `json:"logged_at"`
}

func createUserAndLogin(email, password string) (string, error) {
	user := User{Email: email, Password: password}

	// Create user
	userJSON, _ := json.Marshal(user)
	uri := fmt.Sprintf("http://localhost:%s/api/v1/auth/signup", platform.AppConfig.Port)
	_, err := http.Post(uri, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		return "", fmt.Errorf("signup failed: %v", err)
	}

	// Login and get token
	uri = fmt.Sprintf("http://localhost:%s/api/v1/auth/login", platform.AppConfig.Port)
	resp, err := http.Post(uri, "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		return "", fmt.Errorf("login failed: %v", err)
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	return authResp.Token, nil
}

func addDrinkLogs(token string, baseTime time.Time) error {
	client := &http.Client{}

	// Add three drinks with random spacing (30-90 minutes apart)
	randNumberOfDrinks := 3 + rand.Intn(3) // 3 to 6 drinks

	lastTime := baseTime
	for i := 0; i < randNumberOfDrinks; i++ {
		// Random time interval between 30-90 minutes
		randomMinutes := 30 + rand.Intn(61) // 30 to 90 minutes
		timestamp := lastTime.Add(time.Duration(randomMinutes) * time.Minute)
		lastTime = timestamp

		fmt.Printf("Adding drink for timestamp: %v\n", timestamp)

		drinkLog := DrinkLog{
			DrinkOptionID: rand.Intn(5) + 1, // Random drink option (1-5)
			LoggedAt:      timestamp,
		}

		drinkLogJSON, _ := json.Marshal(drinkLog)
		uri := fmt.Sprintf("http://localhost:%s/api/v1/drink-logs", platform.AppConfig.Port)
		req, err := http.NewRequest("POST", uri, bytes.NewBuffer(drinkLogJSON))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to add drink log: %v", err)
		}
		defer resp.Body.Close()

		// Print response
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Response: %s\n", string(body))

		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {
	baseTime := time.Now().AddDate(0, 0, -1)

	fmt.Println("Creating User 1...")
	token1, err := createUserAndLogin("user1@example.com", "user1-password123")
	fmt.Printf("Token 1: %v\n", token1)
	if err != nil {
		fmt.Printf("Error with user 1: %v\n", err)
		return
	}
	if err := addDrinkLogs(token1, baseTime); err != nil {
		fmt.Printf("Error adding drinks for user 1: %v\n", err)
		return
	}

	fmt.Println("\nCreating User 2...")
	token2, err := createUserAndLogin("user2@example.com", "user2-password123")
	fmt.Printf("Token 2: %v\n", token2)
	if err != nil {
		fmt.Printf("Error with user 2: %v\n", err)
		return
	}
	if err := addDrinkLogs(token2, baseTime.Add(15*time.Minute)); err != nil {
		fmt.Printf("Error adding drinks for user 2: %v\n", err)
		return
	}

	fmt.Println("\nCreating User 3...")
	token3, err := createUserAndLogin("user3@example.com", "user3-password123")
	fmt.Printf("Token 3: %v\n", token3)
	if err != nil {
		fmt.Printf("Error with user 3: %v\n", err)
		return
	}
	if err := addDrinkLogs(token3, baseTime.Add(85*time.Minute)); err != nil {
		fmt.Printf("Error adding drinks for user 3: %v\n", err)
		return
	}
}
