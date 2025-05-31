package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"multi-chain-payment-gateway/internal/config"
	"multi-chain-payment-gateway/internal/models"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

type BlockchainService struct {
	config    *config.Config
	ethClient *ethclient.Client
	wallets   map[string]*WalletInfo
}

type WalletInfo struct {
	Address    string
	PrivateKey string
	Chain      models.Chain
}

func NewBlockchainService(cfg *config.Config) *BlockchainService {
	s := &BlockchainService{
		config:  cfg,
		wallets: make(map[string]*WalletInfo),
	}

	// Initialize Ethereum client if RPC URL is provided
	if cfg.EthereumRPC != "" {
		if client, err := ethclient.Dial(cfg.EthereumRPC); err == nil {
			s.ethClient = client
		}
	}

	return s
}

func (s *BlockchainService) GenerateWallet(chain models.Chain) (*WalletInfo, error) {
	switch chain {
	case models.ChainEthereum:
		return s.generateEthereumWallet()
	case models.ChainSolana:
		return s.generateSolanaWallet()
	case models.ChainTON:
		return s.generateTONWallet()
	default:
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}
}

func (s *BlockchainService) generateEthereumWallet() (*WalletInfo, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	privateKeyHex := fmt.Sprintf("%x", privateKey.D.Bytes())

	wallet := &WalletInfo{
		Address:    address,
		PrivateKey: privateKeyHex,
		Chain:      models.ChainEthereum,
	}

	s.wallets[address] = wallet
	return wallet, nil
}

func (s *BlockchainService) generateSolanaWallet() (*WalletInfo, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("Sol%x", bytes[:16])
	privateKey := fmt.Sprintf("%x", bytes)

	wallet := &WalletInfo{
		Address:    address,
		PrivateKey: privateKey,
		Chain:      models.ChainSolana,
	}

	s.wallets[address] = wallet
	return wallet, nil
}

func (s *BlockchainService) generateTONWallet() (*WalletInfo, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("TON%x", bytes[:16])
	privateKey := fmt.Sprintf("%x", bytes)

	wallet := &WalletInfo{
		Address:    address,
		PrivateKey: privateKey,
		Chain:      models.ChainTON,
	}

	s.wallets[address] = wallet
	return wallet, nil
}

func (s *BlockchainService) CheckTransaction(chain models.Chain, address string, amount decimal.Decimal) (*models.Transaction, error) {
	switch chain {
	case models.ChainEthereum:
		return s.checkEthereumTransaction(address, amount)
	case models.ChainSolana:
		return s.checkSolanaTransaction(address, amount)
	case models.ChainTON:
		return s.checkTONTransaction(address, amount)
	default:
		return nil, fmt.Errorf("unsupported chain: %s", chain)
	}
}

func (s *BlockchainService) checkEthereumTransaction(address string, expectedAmount decimal.Decimal) (*models.Transaction, error) {
	if s.ethClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := s.ethClient.BalanceAt(ctx, common.HexToAddress(address), nil)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *BlockchainService) checkSolanaTransaction(address string, expectedAmount decimal.Decimal) (*models.Transaction, error) {
	return nil, nil
}

func (s *BlockchainService) checkTONTransaction(address string, expectedAmount decimal.Decimal) (*models.Transaction, error) {
	return nil, nil
}

func (s *BlockchainService) GetTokenDecimals(chain models.Chain, token models.TokenType) int {
	switch chain {
	case models.ChainEthereum:
		if token == models.TokenNative {
			return 18
		}
		return 6
	case models.ChainSolana:
		if token == models.TokenNative {
			return 9
		}
		return 6
	case models.ChainTON:
		if token == models.TokenNative {
			return 9
		}
		return 6
	default:
		return 18
	}
}

func (s *BlockchainService) GetTokenSymbol(chain models.Chain, token models.TokenType) string {
	switch chain {
	case models.ChainEthereum:
		switch token {
		case models.TokenNative:
			return "ETH"
		case models.TokenUSDC:
			return "USDC"
		case models.TokenUSDT:
			return "USDT"
		}
	case models.ChainSolana:
		switch token {
		case models.TokenNative:
			return "SOL"
		case models.TokenUSDC:
			return "USDC"
		case models.TokenUSDT:
			return "USDT"
		}
	case models.ChainTON:
		switch token {
		case models.TokenNative:
			return "TON"
		case models.TokenUSDC:
			return "USDC"
		case models.TokenUSDT:
			return "USDT"
		}
	}
	return ""
}
