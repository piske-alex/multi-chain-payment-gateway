package api

import (
	"net/http"
	"multi-chain-payment-gateway/internal/services"

	"github.com/gin-gonic/gin"
)

// BinanceHandler handles Binance API related endpoints
type BinanceHandler struct {
	enhancedPaymentService *services.EnhancedPaymentService
}

// NewBinanceHandler creates a new Binance handler
func NewBinanceHandler(enhancedPaymentService *services.EnhancedPaymentService) *BinanceHandler {
	return &BinanceHandler{
		enhancedPaymentService: enhancedPaymentService,
	}
}

// GetBinanceStatus returns the current status of Binance API integration
func (h *BinanceHandler) GetBinanceStatus(c *gin.Context) {
	status := h.enhancedPaymentService.GetBinanceAPIStatus()
	c.JSON(http.StatusOK, gin.H{
		"binance_api": status,
	})
}

// GetCurrentPrices returns current crypto prices from Binance
func (h *BinanceHandler) GetCurrentPrices(c *gin.Context) {
	prices, err := h.enhancedPaymentService.GetAllPrices()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Failed to fetch prices from Binance API",
			"details": err.Error(),
		})
		return
	}

	// Convert decimal prices to strings for JSON response
	priceStrings := make(map[string]string)
	for symbol, price := range prices {
		priceStrings[symbol] = price.String()
	}

	c.JSON(http.StatusOK, gin.H{
		"prices": priceStrings,
		"source": "binance",
		"timestamp": "now",
	})
}

// GetPaymentWithPriceInfo returns payment details with current price information
func (h *BinanceHandler) GetPaymentWithPriceInfo(c *gin.Context) {
	paymentID := c.Param("id")

	paymentWithPrices, err := h.enhancedPaymentService.GetPaymentWithPriceInfo(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, paymentWithPrices)
}

// GetPriceImpact returns information about price changes since payment creation
func (h *BinanceHandler) GetPriceImpact(c *gin.Context) {
	paymentID := c.Param("id")

	impactInfo, err := h.enhancedPaymentService.GetPriceImpactInfo(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found or price data unavailable"})
		return
	}

	c.JSON(http.StatusOK, impactInfo)
}

// RefreshPaymentPrices updates a payment's crypto amounts with current prices
func (h *BinanceHandler) RefreshPaymentPrices(c *gin.Context) {
	paymentID := c.Param("id")

	err := h.enhancedPaymentService.RefreshPaymentPrices(paymentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to refresh prices",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment prices refreshed successfully",
		"payment_id": paymentID,
	})
}