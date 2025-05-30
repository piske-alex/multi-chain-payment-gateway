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
	apiKey string
	cache  map[string]CachedPrice
}

type CachedPrice struct {
	Price     decimal.Decimal
	ExpiresAt time.Time
}

type CoinGeckoResponse struct {
	Ethereum struct {
		USD decimal.Decimal `json:"usd"`
	} `json:"ethereum"`
	Solana struct {
		USD decimal.Decimal `json:"usd"`
	} `json:"solana"`
	Toncoin struct {
		USD decimal.Decimal `json:"usd"`
	} `json:"the-open-network"`
}

func NewPriceService(apiKey string) *PriceService {
	return &PriceService{
		apiKey: apiKey,
		cache:  make(map[string]CachedPrice),
	}
}

func (s *PriceService) GetPrice(symbol string) (decimal.Decimal, error) {
	// Check cache first
	if cached, exists := s.cache[symbol]; exists && time.Now().Before(cached.ExpiresAt) {
		return cached.Price, nil
	}

	// Fetch from API
	price, err := s.fetchPrice(symbol)
	if err != nil {
		return decimal.Zero, err
	}

	// Cache for 1 minute
	s.cache[symbol] = CachedPrice{
		Price:     price,
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}

	return price, nil
}

func (s *PriceService) fetchPrice(symbol string) (decimal.Decimal, error) {
	var coinId string
	switch symbol {
	case "ETH":
		coinId = "ethereum"
	case "SOL":
		coinId = "solana"
	case "TON":
		coinId = "the-open-network"
	default:
		return decimal.Zero, fmt.Errorf("unsupported symbol: %s", symbol)
	}

	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", coinId)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return decimal.Zero, err
	}

	if s.apiKey != "" {
		req.Header.Set("X-CG-Demo-API-Key", s.apiKey)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}

	var result map[string]map[string]decimal.Decimal
	if err := json.Unmarshal(body, &result); err != nil {
		return decimal.Zero, err
	}

	if priceData, exists := result[coinId]; exists {
		if price, exists := priceData["usd"]; exists {
			return price, nil
		}
	}

	return decimal.Zero, fmt.Errorf("price not found for %s", symbol)
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