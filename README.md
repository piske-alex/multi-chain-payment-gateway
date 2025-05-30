# Multi-Chain Payment Gateway

A comprehensive crypto payment gateway supporting Ethereum, TON, and Solana networks with native tokens and stablecoins (USDC/USDT). Built with Go backend and SvelteKit frontend.

## 🚀 Features

- **Multi-chain Support**: Ethereum, TON, and Solana
- **9 Payment Options**: Native tokens (ETH, TON, SOL) + USDC/USDT on each chain
- **Real-time Price Conversion**: USD to crypto conversion using CoinGecko API
- **Payment Detection**: Monitors blockchain for incoming payments
- **Webhook Integration**: Configurable webhook notifications with HMAC signatures
- **Payment Widget**: Embeddable SvelteKit widget or redirect flow
- **Success Page Redirect**: Configurable success page redirection
- **QR Code Generation**: Built-in QR codes for easy mobile payments

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   SvelteKit     │    │    Go API       │    │   Blockchain    │
│   Frontend      │◄──►│   (Gin + GORM)  │◄──►│   Networks      │
│   Widget        │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │    SQLite DB    │
                       │   (Payments)    │
                       └─────────────────┘
```

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/piske-alex/multi-chain-payment-gateway.git
cd multi-chain-payment-gateway

# Create environment file
cp .env.example .env
# Edit .env with your configuration

# Run with Docker Compose
docker-compose up -d
```

### Manual Setup

#### Backend (Go API)

```bash
# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env
# Configure your environment variables

# Run the server
go run main.go
```

#### Frontend (SvelteKit Widget)

```bash
cd frontend
npm install
npm run dev
```

## 📡 API Endpoints

### Create Payment Intent
```http
POST /api/payments
Content-Type: application/json

{
  "amount": 25.50,
  "currency": "USD",
  "webhook_url": "https://your-site.com/webhook",
  "success_url": "https://your-site.com/success",
  "metadata": {
    "order_id": "12345",
    "customer_id": "user_789"
  }
}
```

**Response:**
```json
{
  "id": "payment_123",
  "amount": "25.50",
  "currency": "USD",
  "status": "pending",
  "expires_at": "2024-01-01T12:30:00Z",
  "options": [
    {
      "chain": "ethereum",
      "token": "native",
      "address": "0x742d35Cc6478354...",
      "amount": "0.01234567",
      "symbol": "ETH",
      "decimals": 18
    }
    // ... 8 more options
  ]
}
```

### Get Payment Details
```http
GET /api/payments/{payment_id}
```

### Get Payment Status
```http
GET /api/payments/{payment_id}/status
```

### Payment Widget
```http
GET /widget/{payment_id}
```

## 🔧 Configuration

### Environment Variables

```bash
# API Configuration
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=sqlite://./payments.db

# Blockchain RPC URLs
ETHEREUM_RPC_URL=https://eth-mainnet.g.alchemy.com/v2/your-key
SOLANA_RPC_URL=https://api.mainnet-beta.solana.com
TON_RPC_URL=https://toncenter.com/api/v2/jsonRPC

# Price API (CoinGecko)
PRICE_API_KEY=your-coingecko-api-key

# Webhook Security
WEBHOOK_SECRET=your-webhook-secret

# Widget Configuration
WIDGET_BASE_URL=http://localhost:5173
```

## 🔗 Integration Examples

### Embed Widget (iframe)
```html
<iframe 
  src="https://your-gateway.com/widget/payment_123" 
  width="400" 
  height="600"
  frameborder="0">
</iframe>
```

### Redirect Integration
```javascript
// Create payment
const response = await fetch('/api/payments', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    amount: 25.50,
    currency: 'USD',
    success_url: 'https://your-site.com/success',
    webhook_url: 'https://your-site.com/webhook'
  })
});

const payment = await response.json();

// Redirect to payment page
window.location.href = `/widget/${payment.id}`;
```

### Webhook Handling
```javascript
// Express.js webhook handler
app.post('/webhook', express.raw({type: 'application/json'}), (req, res) => {
  const signature = req.headers['x-webhook-signature'];
  const payload = req.body;
  
  // Verify signature
  const expectedSignature = crypto
    .createHmac('sha256', process.env.WEBHOOK_SECRET)
    .update(payload)
    .digest('hex');
  
  if (signature === `sha256=${expectedSignature}`) {
    const event = JSON.parse(payload);
    
    if (event.event === 'payment.completed') {
      // Process successful payment
      console.log('Payment completed:', event.payment_id);
    }
  }
  
  res.sendStatus(200);
});
```

## 🔐 Security Features

- **HMAC Webhook Signatures**: All webhooks are signed with HMAC-SHA256
- **Payment Expiration**: Payments automatically expire after 30 minutes
- **Address Generation**: Unique addresses generated for each payment
- **CORS Protection**: Configurable CORS policies
- **Input Validation**: Comprehensive request validation

## 🚦 Payment Flow

1. **Create Payment**: POST to `/api/payments` with amount and metadata
2. **Display Options**: Show 9 payment options (3 chains × 3 tokens)
3. **User Selection**: Customer chooses preferred payment method
4. **Address Display**: Show QR code and wallet address
5. **Monitoring**: System monitors blockchain for incoming transactions
6. **Webhook Notification**: Send webhook when payment is detected
7. **Success Redirect**: Redirect to success URL

## 🔍 Supported Networks

| Network  | Native Token | USDC | USDT |
|----------|--------------|------|------|
| Ethereum | ETH          | ✅    | ✅    |
| Solana   | SOL          | ✅    | ✅    |
| TON      | TON          | ✅    | ✅    |

## 📊 Monitoring & Health

- **Health Check**: `GET /health`
- **Payment Status Polling**: Real-time status updates
- **Webhook Retry Logic**: Automatic retry for failed webhooks
- **Error Handling**: Comprehensive error responses

## 🛠️ Development

### Project Structure
```
├── main.go                 # Application entry point
├── internal/
│   ├── api/               # HTTP handlers and routing
│   ├── config/            # Configuration management
│   ├── database/          # Database initialization
│   ├── models/            # Data models
│   └── services/          # Business logic
├── frontend/              # SvelteKit widget
│   ├── src/
│   │   ├── lib/          # Svelte components
│   │   └── routes/       # Pages
│   └── package.json
├── static/               # Static assets
└── docker-compose.yml    # Docker configuration
```

### Running Tests
```bash
# Backend tests
go test ./...

# Frontend tests
cd frontend
npm test
```

## 📝 API Documentation

See `api-examples.http` for complete API examples and test requests.

## 🚀 Deployment

### Docker
```bash
docker build -t payment-gateway .
docker run -p 8080:8080 payment-gateway
```

### Production Considerations
- Use PostgreSQL instead of SQLite for production
- Configure proper RPC endpoints for each blockchain
- Set up monitoring and alerting
- Use environment-specific secrets
- Enable HTTPS/TLS

## 📄 License

MIT License - see LICENSE file for details.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

For questions and support, please open an issue in the GitHub repository.