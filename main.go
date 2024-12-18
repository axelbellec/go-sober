package main

import (
	"log"
	"net/http"

	"go-sober/internal/analytics"
	"go-sober/internal/auth"
	"go-sober/internal/database"
	"go-sober/internal/drinks"
	"go-sober/internal/middleware"
	"go-sober/platform"
)

func main() {
	// Initialize platform (config, logger, etc)
	platform.InitPlatform()

	// Initialize SQLite database using config
	db, err := database.NewSQLiteDB(platform.AppConfig.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repository, service, controller and middleware with config
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, platform.AppConfig)
	authController := auth.NewController(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize the drinks components
	drinkRepo := drinks.NewRepository(db)
	drinkService := drinks.NewService(drinkRepo)
	drinkController := drinks.NewController(drinkService)

	// Initialize analytics components
	analyticsService := analytics.NewService(drinkRepo)
	analyticsController := analytics.NewController(analyticsService)

	// Public routes
	// Auth
	http.HandleFunc("/auth/signup", authController.SignUp)
	http.HandleFunc("/auth/login", authController.Login)

	// Drink options
	http.HandleFunc("/drink-options", drinkController.GetDrinkOptions)
	http.HandleFunc("/drink-options/", drinkController.GetDrinkOption)

	// Protected routes
	http.HandleFunc("/auth/me", authMiddleware.RequireAuth(authController.Me))
	http.HandleFunc("/analytics/timeline/bac", authMiddleware.RequireAuth(analyticsController.GetBAC))

	// Drink logging
	http.HandleFunc("/drink-logs", authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			drinkController.CreateDrinkLog(w, r)
		case http.MethodGet:
			drinkController.GetDrinkLogs(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))

	// Start server using config port
	addr := ":" + platform.AppConfig.Port
	log.Printf("Server starting on %s, environment: %s\n", addr, platform.AppConfig.Environment)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
