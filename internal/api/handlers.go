package api

import (
	"multi-chain-payment-gateway/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
	webhookService *services.WebhookService
}

func NewPaymentHandler(paymentService *services.PaymentService, webhookService *services.WebhookService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		webhookService: webhookService,
	}
}

func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req services.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate currency
	if req.Currency != "USD" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only USD currency is supported"})
		return
	}

	// Create payment
	payment, err := h.paymentService.CreatePayment(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *PaymentHandler) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentService.GetPayment(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentID := c.Param("id")

	payment, err := h.paymentService.GetPayment(paymentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         payment.ID,
		"status":     payment.Status,
		"amount":     payment.Amount,
		"currency":   payment.Currency,
		"expires_at": payment.ExpiresAt,
	})
}

func (h *PaymentHandler) ServeWidget(c *gin.Context) {
	paymentID := c.Param("id")

	// Check if payment exists
	_, err := h.paymentService.GetPayment(paymentID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"error": "Payment not found",
		})
		return
	}

	// Serve the widget HTML that loads the SvelteKit app
	widgetHTML := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Crypto Payment</title>
    <style>
        body { margin: 0; padding: 20px; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; }
        .widget-container { max-width: 400px; margin: 0 auto; }
    </style>
</head>
<body>
    <div class="widget-container">
        <div id="payment-widget" data-payment-id="` + paymentID + `"></div>
    </div>
    <script>
        window.PAYMENT_ID = '` + paymentID + `';
        window.API_BASE_URL = window.location.origin;
    </script>
    <script src="/static/widget.js"></script>
</body>
</html>`

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, widgetHTML)
}
