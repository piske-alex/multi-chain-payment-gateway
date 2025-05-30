package services

import (
	"multi-chain-payment-gateway/internal/config"
	"multi-chain-payment-gateway/internal/models"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type PaymentService struct {
	db         *gorm.DB
	priceService *PriceService
	blockchainService *BlockchainService
	config     *config.Config
}

type CreatePaymentRequest struct {
	Amount     decimal.Decimal   `json:"amount" binding:"required"`
	Currency   string           `json:"currency" binding:"required"`
	WebhookURL string           `json:"webhook_url"`
	SuccessURL string           `json:"success_url"`
	Metadata   map[string]interface{} `json:"metadata"`
}

func NewPaymentService(db *gorm.DB, priceService *PriceService, blockchainService *BlockchainService, config *config.Config) *PaymentService {
	return &PaymentService{
		db:         db,
		priceService: priceService,
		blockchainService: blockchainService,
		config:     config,
	}
}

func (s *PaymentService) CreatePayment(req CreatePaymentRequest) (*models.Payment, error) {
	// Generate payment ID
	paymentID := uuid.New().String()

	// Serialize metadata
	metadataJSON, _ := json.Marshal(req.Metadata)

	// Create payment
	payment := &models.Payment{
		ID:         paymentID,
		Amount:     req.Amount,
		Currency:   req.Currency,
		Status:     models.StatusPending,
		WebhookURL: req.WebhookURL,
		SuccessURL: req.SuccessURL,
		Metadata:   string(metadataJSON),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}

	// Save payment
	if err := s.db.Create(payment).Error; err != nil {
		return nil, err
	}

	// Generate payment options
	if err := s.generatePaymentOptions(payment); err != nil {
		return nil, err
	}

	// Load payment with options
	if err := s.db.Preload("Options").First(payment, "id = ?", paymentID).Error; err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) generatePaymentOptions(payment *models.Payment) error {
	chains := []models.Chain{models.ChainEthereum, models.ChainSolana, models.ChainTON}
	tokens := []models.TokenType{models.TokenNative, models.TokenUSDC, models.TokenUSDT}

	for _, chain := range chains {
		for _, token := range tokens {
			// Generate wallet for this chain
			wallet, err := s.blockchainService.GenerateWallet(chain)
			if err != nil {
				return err
			}

			// Get token symbol and decimals
			symbol := s.blockchainService.GetTokenSymbol(chain, token)
			decimals := s.blockchainService.GetTokenDecimals(chain, token)

			// Calculate amount in crypto
			var cryptoAmount decimal.Decimal
			if token == models.TokenNative {
				// Convert USD to native token
				var err error
				cryptoAmount, err = s.priceService.ConvertUSDToCrypto(payment.Amount, symbol)
				if err != nil {
					return err
				}
			} else {
				// For stablecoins, amount is 1:1 with USD
				cryptoAmount = payment.Amount
			}

			// Create payment option
			option := &models.PaymentOption{
				PaymentID: payment.ID,
				Chain:     chain,
				Token:     token,
				Address:   wallet.Address,
				Amount:    cryptoAmount,
				Symbol:    symbol,
				Decimals:  decimals,
			}

			if err := s.db.Create(option).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *PaymentService) GetPayment(paymentID string) (*models.Payment, error) {
	var payment models.Payment
	err := s.db.Preload("Options").Preload("Transactions").First(&payment, "id = ?", paymentID).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (s *PaymentService) StartMonitoring() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkPendingPayments()
		}
	}
}

func (s *PaymentService) checkPendingPayments() {
	var payments []models.Payment
	err := s.db.Preload("Options").Where("status = ? AND expires_at > ?", models.StatusPending, time.Now()).Find(&payments).Error
	if err != nil {
		log.Printf("Error fetching pending payments: %v", err)
		return
	}

	for _, payment := range payments {
		for _, option := range payment.Options {
			tx, err := s.blockchainService.CheckTransaction(option.Chain, option.Address, option.Amount)
			if err != nil {
				log.Printf("Error checking transaction for payment %s: %v", payment.ID, err)
				continue
			}

			if tx != nil {
				// Payment received
				s.processPayment(&payment, tx)
				break
			}
		}
	}

	// Mark expired payments
	s.markExpiredPayments()
}

func (s *PaymentService) processPayment(payment *models.Payment, tx *models.Transaction) {
	// Update payment status
	payment.Status = models.StatusPaid
	s.db.Save(payment)

	// Save transaction
	tx.PaymentID = payment.ID
	s.db.Create(tx)

	log.Printf("Payment %s completed with transaction %s", payment.ID, tx.TxHash)
}

func (s *PaymentService) markExpiredPayments() {
	s.db.Model(&models.Payment{}).Where("status = ? AND expires_at <= ?", models.StatusPending, time.Now()).Update("status", models.StatusExpired)
}