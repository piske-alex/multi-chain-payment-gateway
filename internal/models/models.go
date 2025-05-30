package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusPaid      PaymentStatus = "paid"
	StatusExpired   PaymentStatus = "expired"
	StatusCancelled PaymentStatus = "cancelled"
)

type Chain string

const (
	ChainEthereum Chain = "ethereum"
	ChainSolana   Chain = "solana"
	ChainTON      Chain = "ton"
)

type TokenType string

const (
	TokenNative TokenType = "native"
	TokenUSDC   TokenType = "usdc"
	TokenUSDT   TokenType = "usdt"
)

type Payment struct {
	ID          string          `json:"id" gorm:"primaryKey"`
	Amount      decimal.Decimal `json:"amount" gorm:"type:decimal(20,8)"`
	Currency    string          `json:"currency"`
	Status      PaymentStatus   `json:"status"`
	WebhookURL  string          `json:"webhook_url"`
	SuccessURL  string          `json:"success_url"`
	Metadata    string          `json:"metadata" gorm:"type:text"`
	ExpiresAt   time.Time       `json:"expires_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   gorm.DeletedAt  `json:"-" gorm:"index"`

	Options      []PaymentOption `json:"options" gorm:"foreignKey:PaymentID"`
	Transactions []Transaction   `json:"transactions" gorm:"foreignKey:PaymentID"`
}

type PaymentOption struct {
	ID        uint            `json:"id" gorm:"primaryKey"`
	PaymentID string          `json:"payment_id"`
	Chain     Chain           `json:"chain"`
	Token     TokenType       `json:"token"`
	Address   string          `json:"address"`
	Amount    decimal.Decimal `json:"amount" gorm:"type:decimal(20,8)"`
	Symbol    string          `json:"symbol"`
	Decimals  int             `json:"decimals"`
	CreatedAt time.Time       `json:"created_at"`
}

type Transaction struct {
	ID            uint            `json:"id" gorm:"primaryKey"`
	PaymentID     string          `json:"payment_id"`
	Chain         Chain           `json:"chain"`
	TxHash        string          `json:"tx_hash" gorm:"uniqueIndex"`
	FromAddress   string          `json:"from_address"`
	ToAddress     string          `json:"to_address"`
	Amount        decimal.Decimal `json:"amount" gorm:"type:decimal(20,8)"`
	Token         TokenType       `json:"token"`
	Confirmations int             `json:"confirmations"`
	Confirmed     bool            `json:"confirmed"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}