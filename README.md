# Multi-Chain Payment Gateway

A multi-chain crypto payment gateway supporting Ethereum, TON, and Solana networks with native tokens and stablecoins (USDC/USDT).

## Features

- **Multi-chain Support**: Ethereum, TON, and Solana
- **9 Payment Options**: Native tokens (ETH, TON, SOL) + USDC/USDT on each chain
- **Real-time Price Conversion**: USD to crypto conversion using live rates
- **Payment Detection**: Monitors blockchain for incoming payments
- **Webhook Integration**: Configurable webhook notifications
- **Payment Widget**: Embeddable SvelteKit widget or redirect flow
- **Success Page Redirect**: Configurable success page redirection

## Quick Start

### Backend (Go API)

```bash
go mod tidy
cp .env.example .env
# Configure your environment variables
go run main.go
```

### Frontend (SvelteKit Widget)

```bash
cd frontend
npm install
npm run dev
```

## API Endpoints

### Create Payment Intent
```bash
POST /api/payments
{
  "amount": 10.50,
  "currency": "USD",
  "webhook_url": "https://your-site.com/webhook",
  "success_url": "https://your-site.com/success",
  "metadata": {"order_id": "12345"}
}
```

### Get Payment Status
```bash
GET /api/payments/{payment_id}
```

### Payment Widget
```bash
GET /widget/{payment_id}
```

## Environment Variables

```env
# API Configuration
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=sqlite://./payments.db

# Blockchain RPC URLs
ETHEREUM_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/your-key
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
TON_RPC_URL=https://toncenter.com/api/v2/jsonRPC

# Price API
PRICE_API_KEY=your-coingecko-api-key

# Webhook signing
WEBHOOK_SECRET=your-webhook-secret
```

## Architecture

- **Backend**: Go with Gin framework
- **Frontend**: SvelteKit widget
- **Database**: SQLite (configurable)
- **Blockchain**: Multiple RPC providers
- **Price Feed**: CoinGecko API

## License

MIT License