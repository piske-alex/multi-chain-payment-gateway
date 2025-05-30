package api

import (
	"multi-chain-payment-gateway/internal/config"
	"multi-chain-payment-gateway/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(paymentService *services.PaymentService, webhookService *services.WebhookService, cfg *config.Config) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize handlers
	paymentHandler := NewPaymentHandler(paymentService, webhookService)

	// API routes
	api := r.Group("/api")
	{
		api.POST("/payments", paymentHandler.CreatePayment)
		api.GET("/payments/:id", paymentHandler.GetPayment)
		api.GET("/payments/:id/status", paymentHandler.GetPaymentStatus)
	}

	// Widget routes
	r.GET("/widget/:id", paymentHandler.ServeWidget)

	// Static files (for widget assets)
	r.Static("/static", "./static")

	return r
}