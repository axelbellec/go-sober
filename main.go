package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2" // http-swagger middleware

	_ "go-sober/docs" // swagger docs
	"go-sober/internal/analytics"
	"go-sober/internal/auth"
	"go-sober/internal/database"
	"go-sober/internal/drinks"
	"go-sober/internal/embedding"
	"go-sober/internal/middleware"
	"go-sober/platform"
)

// @title Sober API
// @version 1.0
// @description API for the Sober app
// @host localhost:8080
// @BasePath /api/v1
// @accept json
// @produce json
func main() {
	// Initialize platform (config, logger, etc)
	platform.InitPlatform()

	config := platform.AppConfig

	loggingMiddleware := middleware.NewLoggingMiddleware()

	// Initialize SQLite database using config
	db, err := database.NewSQLiteDB(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize repository, service, controller and middleware with config
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, config)
	authController := auth.NewController(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)

	// Initialize the drinks components
	drinkRepo := drinks.NewRepository(db)
	drinkService := drinks.NewService(drinkRepo)
	embeddingService := embedding.NewOllamaEmbedding(
		config.Embedding.Ollama.BaseURL,
		config.Embedding.Ollama.Model,
	)
	drinkController := drinks.NewController(drinkService, embeddingService, db)

	// Initialize analytics components
	analyticsService := analytics.NewService(drinkRepo)
	analyticsController := analytics.NewController(analyticsService)

	// Create a new ServeMux to use with the logging middleware
	mux := http.NewServeMux()

	// [Public routes]
	// Auth
	mux.HandleFunc("POST /api/v1/auth/signup", authController.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", authController.Login)

	// Drink options
	mux.HandleFunc("GET /api/v1/drink-options", drinkController.GetDrinkOptions)
	mux.HandleFunc("GET /api/v1/drink-options/{id}", drinkController.GetDrinkOption)
	mux.HandleFunc("PUT /api/v1/drink-options/{id}", drinkController.UpdateDrinkOption)
	mux.HandleFunc("DELETE /api/v1/drink-options/{id}", drinkController.DeleteDrinkOption)

	// [Protected routes]
	// Auth
	mux.HandleFunc("GET /api/v1/auth/me", authMiddleware.RequireAuth(authController.Me))

	// Analytics
	mux.HandleFunc("GET /api/v1/analytics/timeline/bac", authMiddleware.RequireAuth(analyticsController.GetBAC))

	// Drink logging
	mux.HandleFunc("GET /api/v1/drink-logs", authMiddleware.RequireAuth(drinkController.GetDrinkLogs))
	mux.HandleFunc("POST /api/v1/drink-logs", authMiddleware.RequireAuth(drinkController.CreateDrinkLog))
	mux.HandleFunc("POST /api/v1/drink-logs/parse", authMiddleware.RequireAuth(drinkController.ParseDrinkLog))

	// Swagger documentation
	mux.HandleFunc("GET /api/v1/swagger/doc.json", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /api/v1/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/api/v1/swagger/doc.json"), // Use relative URL
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Start server using config port
	addr := ":" + config.Port
	log.Printf("Server starting on http://localhost%s, environment: %s\n", addr, config.Environment)
	if err := http.ListenAndServe(addr, loggingMiddleware.LogRequest(mux)); err != nil {
		log.Fatal(err)
	}

}
