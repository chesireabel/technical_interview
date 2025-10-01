package main

import (
	"context"
	"log"
	"os"

	"github.com/chesireabel/Technical-Interview/config"
	"github.com/chesireabel/Technical-Interview/database"
	"github.com/chesireabel/Technical-Interview/internal/handlers"
	"github.com/chesireabel/Technical-Interview/internal/repositories"
	"github.com/chesireabel/Technical-Interview/internal/routes"
	"github.com/chesireabel/Technical-Interview/internal/services"
	"github.com/chesireabel/Technical-Interview/internal/middleware"


	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// Connect to database
	database.ConnectDB()
	defer database.CloseDB()

	// Run migrations
	log.Println("Running database migrations...")
	if err := RunMigrations(database.DB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

     oidc, err := config.InitOIDCWithDefaults(context.Background())
	if err != nil {
		log.Fatalf("❌ Failed to initialize OIDC: %v", err)
	}

	returnToURL := config.GetReturnURL()

	// Initialize SMS service
	smsService, err := services.NewSMSService()
	if err != nil {
		log.Printf("⚠️ Warning: SMS service initialization failed: %v", err)
		log.Println("ℹ️ Orders will be created without SMS notifications")
		smsService = nil // Continue without SMS
	} else {
		log.Println("✅ SMS service initialized successfully")
	}

	// Initialize repositories
	customerRepo := repositories.NewCustomerRepository(database.DB)
	orderRepo := repositories.NewOrderRepository(database.DB)

	// Initialize services
	customerService := services.NewCustomerService(customerRepo)
	orderService := services.NewOrderService(orderRepo, customerRepo, smsService)

	// Initialize handlers
	customerHandler := handlers.NewCustomerHandler(customerService)
	orderHandler := handlers.NewOrderHandler(orderService)

	// Setup a Gin router
	r := gin.Default()

	r.Use(middleware.InitSessionStore())

	// Register all routes
	routes.RegisterRoutes(r, customerHandler, orderHandler, oidc ,returnToURL)

	// Get port from .env
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server is running on: http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}