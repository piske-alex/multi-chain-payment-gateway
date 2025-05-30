package config

import "os"

type Config struct {
	Environment    string
	Port          string
	DatabaseURL   string
	EthereumRPC   string
	SolanaRPC     string
	TonRPC        string
	PriceAPIKey   string
	WebhookSecret string
	WidgetBaseURL string
}

func Load() *Config {
	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "sqlite://./payments.db"),
		EthereumRPC:   getEnv("ETHEREUM_RPC_URL", ""),
		SolanaRPC:     getEnv("SOLANA_RPC_URL", "https://api.mainnet-beta.solana.com"),
		TonRPC:        getEnv("TON_RPC_URL", "https://toncenter.com/api/v2/jsonRPC"),
		PriceAPIKey:   getEnv("PRICE_API_KEY", ""),
		WebhookSecret: getEnv("WEBHOOK_SECRET", "default-secret"),
		WidgetBaseURL: getEnv("WIDGET_BASE_URL", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}