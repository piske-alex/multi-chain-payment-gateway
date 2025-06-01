package services

import (
	"encoding/json"
	"log"
	"time"

	"multi-chain-payment-gateway/internal/models"
	"github.com/shopspring/decimal"
)

// EnhancedPaymentService provides additional methods for better payment handling with Binance integration
type EnhancedPaymentService struct {
	*PaymentService
}

// NewEnhancedPaymentService creates an enhanced payment service
func NewEnhancedPaymentService(paymentService *PaymentService) *EnhancedPaymentService {
	return &EnhancedPaymentService{
		PaymentService: paymentService,
	}
}

// CreatePaymentWithPriceValidation creates a payment with additional price validation
func (s *EnhancedPaymentService) CreatePaymentWithPriceValidation(req CreatePaymentRequest) (*models.Payment, error) {
	// Pre-validate that we can fetch prices for all required tokens
	prices, err := s.priceService.GetAllPrices()
	if err != nil {
		log.Printf("Warning: Could not pre-fetch prices from Binance: %v", err)
		// Continue anyway, individual price fetching might work
	} else {
		log.Printf("Successfully pre-fetched prices: ETH=$%s, SOL=$%s, TON=$%s", 
			prices["ETH"].String(), prices["SOL"].String(), prices["TON"].String())
	}

	return s.PaymentService.CreatePayment(req)
}

// GetPaymentWithPriceInfo returns payment with current price information
func (s *EnhancedPaymentService) GetPaymentWithPriceInfo(paymentID string) (*PaymentWithPriceInfo, error) {
	payment, err := s.PaymentService.GetPayment(paymentID)
	if err != nil {
		return nil, err
	}

	// Get current prices for comparison
	currentPrices, err := s.priceService.GetAllPrices()
	if err != nil {
		log.Printf("Could not fetch current prices: %v", err)
		currentPrices = make(map[string]decimal.Decimal)
	}

	return &PaymentWithPriceInfo{
		Payment:       payment,
		CurrentPrices: currentPrices,
		PriceAge:      time.Now(),
	}, nil
}

// PaymentWithPriceInfo extends payment with current price data
type PaymentWithPriceInfo struct {
	*models.Payment
	CurrentPrices map[string]decimal.Decimal `json:"current_prices"`
	PriceAge      time.Time                  `json:"price_age"`
}

// RefreshPaymentPrices updates payment options with current prices
func (s *EnhancedPaymentService) RefreshPaymentPrices(paymentID string) error {
	payment, err := s.PaymentService.GetPayment(paymentID)
	if err != nil {
		return err
	}

	// Only refresh if payment is still pending
	if payment.Status != models.StatusPending {
		return nil
	}

	// Get current prices
	currentPrices, err := s.priceService.GetAllPrices()
	if err != nil {
		return err
	}

	// Update each native token option with current price
	for _, option := range payment.Options {
		if option.Token == models.TokenNative {
			if currentPrice, exists := currentPrices[option.Symbol]; exists {
				// Recalculate crypto amount with current price
				newAmount := payment.Amount.Div(currentPrice)
				
				// Update the option in database
				s.db.Model(&option).Update("amount", newAmount)
				log.Printf("Updated %s amount for payment %s: %s (price: $%s)", 
					option.Symbol, paymentID, newAmount.String(), currentPrice.String())
			}
		}
	}

	return nil
}

// GetPriceImpactInfo returns information about price changes since payment creation
func (s *EnhancedPaymentService) GetPriceImpactInfo(paymentID string) (*PriceImpactInfo, error) {
	payment, err := s.PaymentService.GetPayment(paymentID)
	if err != nil {
		return nil, err
	}

	currentPrices, err := s.priceService.GetAllPrices()
	if err != nil {
		return nil, err
	}

	impactInfo := &PriceImpactInfo{
		PaymentID:     paymentID,
		CreatedAt:     payment.CreatedAt,
		CurrentPrices: currentPrices,
		PriceChanges:  make(map[string]PriceChange),
	}

	// Calculate price changes for each native token
	for _, option := range payment.Options {
		if option.Token == models.TokenNative {
			if currentPrice, exists := currentPrices[option.Symbol]; exists {
				// Calculate original price from payment amount and crypto amount
				originalPrice := payment.Amount.Div(option.Amount)
				priceChange := currentPrice.Sub(originalPrice).Div(originalPrice).Mul(decimal.NewFromInt(100))
				
				impactInfo.PriceChanges[option.Symbol] = PriceChange{
					OriginalPrice: originalPrice,
					CurrentPrice:  currentPrice,
					ChangePercent: priceChange,
				}
			}
		}
	}

	return impactInfo, nil
}

type PriceImpactInfo struct {
	PaymentID     string                     `json:"payment_id"`
	CreatedAt     time.Time                  `json:"created_at"`
	CurrentPrices map[string]decimal.Decimal `json:"current_prices"`
	PriceChanges  map[string]PriceChange     `json:"price_changes"`
}

type PriceChange struct {
	OriginalPrice decimal.Decimal `json:"original_price"`
	CurrentPrice  decimal.Decimal `json:"current_price"`
	ChangePercent decimal.Decimal `json:"change_percent"`
}

// ValidateBinanceConnectivity tests the connection to Binance API
func (s *EnhancedPaymentService) ValidateBinanceConnectivity() error {
	_, err := s.priceService.GetPrice("ETH")
	if err != nil {
		return err
	}
	log.Println("Binance API connectivity validated successfully")
	return nil
}

// GetBinanceAPIStatus returns the current status of Binance API
func (s *EnhancedPaymentService) GetBinanceAPIStatus() map[string]interface{} {
	status := make(map[string]interface{})
	
	// Test connectivity
	start := time.Now()
	prices, err := s.priceService.GetAllPrices()
	latency := time.Since(start)
	
	if err != nil {
		status["status"] = "error"
		status["error"] = err.Error()
		status["latency_ms"] = latency.Milliseconds()
	} else {
		status["status"] = "healthy"
		status["latency_ms"] = latency.Milliseconds()
		status["prices_available"] = len(prices)
		status["last_update"] = time.Now().Format(time.RFC3339)
		
		// Include current prices
		priceStrings := make(map[string]string)
		for symbol, price := range prices {
			priceStrings[symbol] = price.String()
		}
		status["prices"] = priceStrings
	}
	
	return status
}