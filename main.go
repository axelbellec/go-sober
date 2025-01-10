package main

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2" // http-swagger middleware

	_ "go-sober/docs" // swagger docs
	"go-sober/internal/analytics"
	"go-sober/internal/auth"
	"go-sober/internal/bac"
	"go-sober/internal/database"
	"go-sober/internal/drinks"
	"go-sober/internal/embedding"
	"go-sober/internal/health"
	"go-sober/internal/middleware"
	"go-sober/internal/user"
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

	corsMiddleware := middleware.NewCorsMiddleware()
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

	// Initialize the health components
	healthController := health.NewController()

	// Initialize the drinks components
	drinkRepo := drinks.NewRepository(db)
	drinkService := drinks.NewService(drinkRepo)
	embeddingService := embedding.NewOllamaEmbedding(
		config.Embedding.Ollama.BaseURL,
		config.Embedding.Ollama.Model,
	)
	drinkController := drinks.NewController(drinkService, embeddingService, db)

	// Initialize analytics components
	drinkStatsRepo := analytics.NewRepository(db)
	drinkStatsService := analytics.NewService(drinkStatsRepo)
	drinkStatsController := analytics.NewController(drinkStatsService)

	// Initialize realtime components
	bacService := bac.NewService(drinkRepo)
	bacController := bac.NewController(bacService)

	// Initialize user components
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userController := user.NewController(userService)

	// Create a new ServeMux to use with the logging middleware
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/health", healthController.Health)

	// [Public routes]
	// Auth
	mux.HandleFunc("POST /api/v1/auth/signup", authController.SignUp)
	mux.HandleFunc("POST /api/v1/auth/login", authController.Login)

	// User
	mux.HandleFunc("GET /api/v1/users/profile", authMiddleware.RequireAuth(userController.GetProfile))
	mux.HandleFunc("PUT /api/v1/users/profile", authMiddleware.RequireAuth(userController.UpdateProfile))

	// Drink templates
	mux.HandleFunc("GET /api/v1/drink-templates", drinkController.GetDrinkTemplates)
	mux.HandleFunc("GET /api/v1/drink-templates/{id}", drinkController.GetDrinkTemplate)
	mux.HandleFunc("PUT /api/v1/drink-templates/{id}", drinkController.UpdateDrinkTemplate)
	mux.HandleFunc("DELETE /api/v1/drink-templates/{id}", drinkController.DeleteDrinkTemplate)

	// [Protected routes]
	// Auth
	mux.HandleFunc("GET /api/v1/auth/me", authMiddleware.RequireAuth(authController.Me))

	// Blood Alcohol Content (BAC)
	mux.HandleFunc("GET /api/v1/bac/timeline", authMiddleware.RequireAuth(bacController.GetBAC))
	mux.HandleFunc("GET /api/v1/bac/current", authMiddleware.RequireAuth(bacController.GetCurrentBAC))

	// Drink logging
	mux.HandleFunc("GET /api/v1/drink-logs", authMiddleware.RequireAuth(drinkController.GetDrinkLogs))
	mux.HandleFunc("POST /api/v1/drink-logs", authMiddleware.RequireAuth(drinkController.CreateDrinkLog))
	mux.HandleFunc("PUT /api/v1/drink-logs", authMiddleware.RequireAuth(drinkController.UpdateDrinkLog))
	mux.HandleFunc("DELETE /api/v1/drink-logs/{id}", authMiddleware.RequireAuth(drinkController.DeleteDrinkLog))
	mux.HandleFunc("POST /api/v1/drink-logs/parse", authMiddleware.RequireAuth(drinkController.ParseDrinkLog))

	// Analytics
	mux.HandleFunc("GET /api/v1/analytics/drink-stats", authMiddleware.RequireAuth(drinkStatsController.GetDrinkStats))
	mux.HandleFunc("GET /api/v1/analytics/monthly-bac", authMiddleware.RequireAuth(drinkStatsController.GetMonthlyBACStats))

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
	if err := http.ListenAndServe(addr, loggingMiddleware.LogRequest(corsMiddleware.EnableCors(mux))); err != nil {
		log.Fatal(err)
	}

}
