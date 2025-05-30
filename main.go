package main

import (
	"log"

	"multi-chain-payment-gateway/internal/api"
	"multi-chain-payment-gateway/internal/config"
	"multi-chain-payment-gateway/internal/database"
	"multi-chain-payment-gateway/internal/services"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	priceService := services.NewPriceService(cfg.PriceAPIKey)
	blockchainService := services.NewBlockchainService(cfg)
	paymentService := services.NewPaymentService(db, priceService, blockchainService, cfg)
	webhookService := services.NewWebhookService(cfg.WebhookSecret)

	// Start blockchain monitoring
	go paymentService.StartMonitoring()

	// Initialize API server
	router := api.NewRouter(paymentService, webhookService, cfg)

	// Start server
	port := cfg.Port

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
