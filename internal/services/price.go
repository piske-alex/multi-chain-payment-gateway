package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
)

type PriceService struct {
	cache  map[string]CachedPrice
}

type CachedPrice struct {
	Price     decimal.Decimal
	ExpiresAt time.Time
}

type BinanceTickerResponse struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
}

type BinancePriceResponse []BinanceTickerResponse

func NewPriceService(apiKey string) *PriceService {
	// Note: Binance public API doesn't require API key for price data
	return &PriceService{
		cache: make(map[string]CachedPrice),
	}
}

func (s *PriceService) GetPrice(symbol string) (decimal.Decimal, error) {
	// Check cache first
	if cached, exists := s.cache[symbol]; exists && time.Now().Before(cached.ExpiresAt) {
		return cached.Price, nil
	}

	// Fetch from Binance API
	price, err := s.fetchPrice(symbol)
	if err != nil {
		return decimal.Zero, err
	}

	// Cache for 30 seconds (Binance updates frequently)
	s.cache[symbol] = CachedPrice{
		Price:     price,
		ExpiresAt: time.Now().Add(30 * time.Second),
	}

	return price, nil
}

func (s *PriceService) fetchPrice(symbol string) (decimal.Decimal, error) {
	// Map crypto symbols to Binance trading pairs
	var binanceSymbol string
	switch symbol {
	case "ETH":
		binanceSymbol = "ETHUSDT"
	case "SOL":
		binanceSymbol = "SOLUSDT"
	case "TON":
		binanceSymbol = "TONUSDT"
	default:
		return decimal.Zero, fmt.Errorf("unsupported symbol: %s", symbol)
	}

	// Use Binance API v3 ticker price endpoint
	url := fmt.Sprintf("https://api.binance.com/api/v3/ticker/price?symbol=%s", binanceSymbol)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return decimal.Zero, err
	}

	// Set User-Agent to identify our application
	req.Header.Set("User-Agent", "MultiChainPaymentGateway/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return decimal.Zero, fmt.Errorf("binance API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}

	var ticker BinanceTickerResponse
	if err := json.Unmarshal(body, &ticker); err != nil {
		return decimal.Zero, err
	}

	if ticker.Price.IsZero() {
		return decimal.Zero, fmt.Errorf("invalid price received for %s", symbol)
	}

	return ticker.Price, nil
}

// GetAllPrices fetches all supported crypto prices in a single request for efficiency
func (s *PriceService) GetAllPrices() (map[string]decimal.Decimal, error) {
	// Fetch multiple symbols at once using Binance batch endpoint
	symbols := []string{"ETHUSDT", "SOLUSDT", "TONUSDT"}
	url := "https://api.binance.com/api/v3/ticker/price"
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "MultiChainPaymentGateway/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("binance API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tickers BinancePriceResponse
	if err := json.Unmarshal(body, &tickers); err != nil {
		return nil, err
	}

	// Map Binance symbols back to our crypto symbols
	prices := make(map[string]decimal.Decimal)
	for _, ticker := range tickers {
		switch ticker.Symbol {
		case "ETHUSDT":
			prices["ETH"] = ticker.Price
			// Cache the price
			s.cache["ETH"] = CachedPrice{
				Price:     ticker.Price,
				ExpiresAt: time.Now().Add(30 * time.Second),
			}
		case "SOLUSDT":
			prices["SOL"] = ticker.Price
			s.cache["SOL"] = CachedPrice{
				Price:     ticker.Price,
				ExpiresAt: time.Now().Add(30 * time.Second),
			}
		case "TONUSDT":
			prices["TON"] = ticker.Price
			s.cache["TON"] = CachedPrice{
				Price:     ticker.Price,
				ExpiresAt: time.Now().Add(30 * time.Second),
			}
		}
	}

	return prices, nil
}

func (s *PriceService) ConvertUSDToCrypto(usdAmount decimal.Decimal, symbol string) (decimal.Decimal, error) {
	price, err := s.GetPrice(symbol)
	if err != nil {
		return decimal.Zero, err
	}

	if price.IsZero() {
		return decimal.Zero, fmt.Errorf("invalid price for %s", symbol)
	}

	return usdAmount.Div(price), nil
}

// GetPriceWithFallback attempts to get price with retry logic
func (s *PriceService) GetPriceWithFallback(symbol string) (decimal.Decimal, error) {
	var lastErr error
	
	// Try up to 3 times with exponential backoff
	for attempt := 1; attempt <= 3; attempt++ {
		price, err := s.GetPrice(symbol)
		if err == nil {
			return price, nil
		}
		
		lastErr = err
		if attempt < 3 {
			// Exponential backoff: 1s, 2s, 4s
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	
	return decimal.Zero, fmt.Errorf("failed to get price after 3 attempts: %w", lastErr)
}