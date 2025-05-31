package config

import "os"

type Config struct {
	Environment      string
	Port            string
	DatabaseURL     string
	EthereumRPC     string
	SolanaRPC       string
	TonRPC          string
	BinanceAPIKey   string
	WebhookSecret   string
	WidgetBaseURL   string
	BinanceRateLimit int
	BinanceRateLimitWindow int
}

func Load() *Config {
	return &Config{
		Environment:      getEnv("ENVIRONMENT", "development"),
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", "sqlite://./payments.db"),
		EthereumRPC:     getEnv("ETHEREUM_RPC_URL", ""),
		SolanaRPC:       getEnv("SOLANA_RPC_URL", "https://api.mainnet-beta.solana.com"),
		TonRPC:          getEnv("TON_RPC_URL", "https://toncenter.com/api/v2/jsonRPC"),
		BinanceAPIKey:   getEnv("BINANCE_API_KEY", ""), // Optional for higher rate limits
		WebhookSecret:   getEnv("WEBHOOK_SECRET", "default-secret"),
		WidgetBaseURL:   getEnv("WIDGET_BASE_URL", "http://localhost:5173"),
		BinanceRateLimit: getEnvInt("BINANCE_RATE_LIMIT", 1200), // requests per window
		BinanceRateLimitWindow: getEnvInt("BINANCE_RATE_LIMIT_WINDOW", 60), // seconds
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue := parseInt(value); intValue > 0 {
			return intValue
		}
	}
	return defaultValue
}

func parseInt(s string) int {
	result := 0
	for _, char := range s {
		if char >= '0' && char <= '9' {
			result = result*10 + int(char-'0')
		} else {
			return 0
		}
	}
	return result
}