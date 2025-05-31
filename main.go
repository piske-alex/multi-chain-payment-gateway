package main

import (
	"multi-chain-payment-gateway/internal/api"
	"multi-chain-payment-gateway/internal/config"
	"multi-chain-payment-gateway/internal/database"
	"multi-chain-payment-gateway/internal/services"
	"log"
	"os"

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
	priceService := services.NewPriceService(cfg.BinanceAPIKey) // Updated to use Binance
	blockchainService := services.NewBlockchainService(cfg)
	paymentService := services.NewPaymentService(db, priceService, blockchainService, cfg)
	webhookService := services.NewWebhookService(cfg.WebhookSecret)

	// Start blockchain monitoring
	go paymentService.StartMonitoring()

	// Initialize API server
	router := api.NewRouter(paymentService, webhookService, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("Using Binance API for crypto price data")
	if cfg.BinanceAPIKey != "" {
		log.Printf("Binance API key configured for higher rate limits")
	} else {
		log.Printf("Using Binance public API (no API key required)")
	}
	
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}