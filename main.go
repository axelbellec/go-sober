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

	loggingMiddleware := middleware.NewLoggingMiddleware()

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

	// Create a new ServeMux to use with the logging middleware
	mux := http.NewServeMux()

	// Public routes
	// Auth
	mux.HandleFunc("/auth/signup", authController.SignUp)
	mux.HandleFunc("/auth/login", authController.Login)

	// Drink options
	mux.HandleFunc("/drink-options", drinkController.GetDrinkOptions)
	mux.HandleFunc("/drink-options/", drinkController.GetDrinkOption)

	// Protected routes
	mux.HandleFunc("/auth/me", authMiddleware.RequireAuth(authController.Me))
	mux.HandleFunc("/analytics/timeline/bac", authMiddleware.RequireAuth(analyticsController.GetBAC))

	// Drink logging
	mux.HandleFunc("/drink-logs", authMiddleware.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
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
	if err := http.ListenAndServe(addr, loggingMiddleware.LogRequest(mux)); err != nil {
		log.Fatal(err)
	}

}
