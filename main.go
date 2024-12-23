package main

import (
	"log"
	"net/http"

	"go-sober/internal/analytics"
	"go-sober/internal/auth"
	"go-sober/internal/database"
	"go-sober/internal/drinks"
	"go-sober/internal/embedding"
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
	embeddingService := embedding.NewOllamaEmbedding("http://localhost:11434")
	drinkController := drinks.NewController(drinkService, embeddingService, db)

	// Initialize analytics components
	analyticsService := analytics.NewService(drinkRepo)
	analyticsController := analytics.NewController(analyticsService)

	// Create a new ServeMux to use with the logging middleware
	mux := http.NewServeMux()

	// [Public routes]
	// Auth
	mux.HandleFunc("POST /auth/signup", authController.SignUp)
	mux.HandleFunc("POST /auth/login", authController.Login)

	// Drink options
	mux.HandleFunc("GET /drink-options", drinkController.GetDrinkOptions)
	mux.HandleFunc("GET /drink-options/", drinkController.GetDrinkOption)
	mux.HandleFunc("PUT /drink-options/", drinkController.UpdateDrinkOption)
	mux.HandleFunc("DELETE /drink-options/", drinkController.DeleteDrinkOption)

	// [Protected routes]
	// Auth
	mux.HandleFunc("GET /auth/me", authMiddleware.RequireAuth(authController.Me))

	// Analytics
	mux.HandleFunc("GET /analytics/timeline/bac", authMiddleware.RequireAuth(analyticsController.GetBAC))

	// Drink logging
	mux.HandleFunc("GET /drink-logs", authMiddleware.RequireAuth(drinkController.GetDrinkLogs))
	mux.HandleFunc("POST /drink-logs", authMiddleware.RequireAuth(drinkController.CreateDrinkLog))
	mux.HandleFunc("POST /drink-logs/parse", authMiddleware.RequireAuth(drinkController.ParseDrinkLog))

	// Start server using config port
	addr := ":" + platform.AppConfig.Port
	log.Printf("Server starting on %s, environment: %s\n", addr, platform.AppConfig.Environment)
	if err := http.ListenAndServe(addr, loggingMiddleware.LogRequest(mux)); err != nil {
		log.Fatal(err)
	}

}
