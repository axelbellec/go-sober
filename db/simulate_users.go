package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"go-sober/internal/models"
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

// Définition des boissons possibles
var drinks = []models.DrinkLog{
	{Name: "Beer", Type: "beer", SizeValue: 500, SizeUnit: "ml", ABV: 0.055},
	{Name: "Wine", Type: "wine", SizeValue: 150, SizeUnit: "ml", ABV: 0.13},
	{Name: "Whiskey", Type: "spirit", SizeValue: 45, SizeUnit: "ml", ABV: 0.40},
	{Name: "Cocktail", Type: "cocktail", SizeValue: 200, SizeUnit: "ml", ABV: 0.12},
	{Name: "Strong Beer", Type: "beer", SizeValue: 330, SizeUnit: "ml", ABV: 0.085},
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

func addDrinkLogsForYear(token string, startDate time.Time) error {
	client := &http.Client{}

	for i := 0; i < 365; i++ {
		// 70% chance de boire ce jour
		if rand.Float32() > 0.7 {
			continue
		}

		baseTime := startDate.AddDate(0, 0, -i)
		// Commence à boire entre 17h et 21h
		baseTime = time.Date(
			baseTime.Year(), baseTime.Month(), baseTime.Day(),
			17+rand.Intn(4), // Heure entre 17-21
			rand.Intn(60),   // Minute aléatoire
			0, 0,
			baseTime.Location(),
		)

		// Nombre aléatoire de verres (2-6)
		numDrinks := 2 + rand.Intn(5)
		lastTime := baseTime

		for d := 0; d < numDrinks; d++ {
			// 30-90 minutes entre les verres
			randomMinutes := 30 + rand.Intn(61)
			timestamp := lastTime.Add(time.Duration(randomMinutes) * time.Minute)

			// Sélectionne une boisson aléatoire
			drink := drinks[rand.Intn(len(drinks))]
			drink.LoggedAt = timestamp.UTC()

			drinkLogJSON, _ := json.Marshal(drink)
			uri := "http://localhost:8080/api/v1/drink-logs"
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
			resp.Body.Close()
		}
	}
	return nil
}

func main() {
	startDate := time.Now()

	fmt.Println("Creating User 1 with year-long history...")
	token1, err := createUserAndLogin("user1@example.com", "user1-password123")
	if err != nil {
		fmt.Printf("Error with user 1: %v\n", err)
		return
	}
	if err := addDrinkLogsForYear(token1, startDate); err != nil {
		fmt.Printf("Error adding drinks for user 1: %v\n", err)
		return
	}
}
