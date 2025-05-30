// Fallback widget implementation for non-SvelteKit environments
// This provides basic functionality when the full SvelteKit app isn't available

(function() {
    const paymentId = window.PAYMENT_ID;
    const apiBaseUrl = window.API_BASE_URL || window.location.origin;
    
    if (!paymentId) {
        console.error('Payment ID not found');
        return;
    }

    const container = document.getElementById('payment-widget');
    if (!container) {
        console.error('Widget container not found');
        return;
    }

    let selectedOption = null;
    let payment = null;

    function createWidget() {
        container.innerHTML = `
            <div style="border: 1px solid #e5e7eb; border-radius: 8px; padding: 24px; background: white; max-width: 400px;">
                <h3 style="margin: 0 0 16px 0; color: #111827; font-size: 20px; font-weight: 600;">Complete Your Payment</h3>
                <div id="loading" style="text-align: center; padding: 40px;">
                    <div style="display: inline-block; width: 32px; height: 32px; border: 3px solid #f3f4f6; border-top: 3px solid #3b82f6; border-radius: 50%; animation: spin 1s linear infinite;"></div>
                    <p style="margin: 16px 0 0 0; color: #6b7280;">Loading payment options...</p>
                </div>
                <div id="payment-options" style="display: none;"></div>
                <div id="payment-details" style="display: none;"></div>
            </div>
            <style>
                @keyframes spin {
                    0% { transform: rotate(0deg); }
                    100% { transform: rotate(360deg); }
                }
                .payment-option {
                    border: 1px solid #e5e7eb;
                    border-radius: 8px;
                    padding: 16px;
                    margin-bottom: 8px;
                    cursor: pointer;
                    transition: all 0.2s;
                }
                .payment-option:hover {
                    border-color: #3b82f6;
                    background-color: #eff6ff;
                }
                .chain-badge {
                    display: inline-block;
                    padding: 2px 8px;
                    border-radius: 12px;
                    font-size: 12px;
                    font-weight: 500;
                }
                .chain-ethereum { background: #dbeafe; color: #1e40af; }
                .chain-solana { background: #e9d5ff; color: #7c3aed; }
                .chain-ton { background: #d1fae5; color: #059669; }
            </style>
        `;
    }

    async function loadPayment() {
        try {
            const response = await fetch(`${apiBaseUrl}/api/payments/${paymentId}`);
            payment = await response.json();
            
            if (payment.status === 'paid') {
                showSuccess();
                return;
            }
            
            if (payment.status === 'expired') {
                showExpired();
                return;
            }
            
            showPaymentOptions();
            startStatusPolling();
        } catch (error) {
            console.error('Error loading payment:', error);
            showError('Failed to load payment information');
        }
    }

    function showPaymentOptions() {
        const loading = document.getElementById('loading');
        const optionsContainer = document.getElementById('payment-options');
        
        loading.style.display = 'none';
        optionsContainer.style.display = 'block';
        
        const chains = {
            'ethereum': 'Ethereum',
            'solana': 'Solana', 
            'ton': 'TON'
        };
        
        let html = `
            <div style="background: #f9fafb; border-radius: 8px; padding: 16px; margin-bottom: 20px;">
                <div style="display: flex; justify-content: space-between; align-items: center;">
                    <span style="color: #6b7280; font-size: 14px;">Amount:</span>
                    <span style="font-size: 18px; font-weight: 600;">$${payment.amount} ${payment.currency}</span>
                </div>
            </div>
            <h4 style="margin: 0 0 12px 0; color: #374151; font-size: 14px; font-weight: 500;">Choose Payment Method:</h4>
        `;
        
        payment.options.forEach(option => {
            const chainName = chains[option.chain] || option.chain;
            const amount = parseFloat(option.amount).toFixed(Math.min(option.decimals, 8));
            
            html += `
                <div class="payment-option" onclick="selectOption(${option.id})">
                    <div style="display: flex; justify-content: space-between; align-items: center;">
                        <div>
                            <div style="display: flex; align-items: center; margin-bottom: 4px;">
                                <span style="font-weight: 500; margin-right: 8px;">${option.symbol}</span>
                                <span class="chain-badge chain-${option.chain}">${chainName}</span>
                            </div>
                            <div style="color: #6b7280; font-size: 14px;">${amount} ${option.symbol}</div>
                        </div>
                        <span style="color: #9ca3af;">→</span>
                    </div>
                </div>
            `;
        });
        
        optionsContainer.innerHTML = html;
    }

    window.selectOption = function(optionId) {
        selectedOption = payment.options.find(opt => opt.id === optionId);
        showPaymentDetails();
    }

    function showPaymentDetails() {
        const optionsContainer = document.getElementById('payment-options');
        const detailsContainer = document.getElementById('payment-details');
        
        optionsContainer.style.display = 'none';
        detailsContainer.style.display = 'block';
        
        const chains = {
            'ethereum': 'Ethereum',
            'solana': 'Solana', 
            'ton': 'TON'
        };
        
        const chainName = chains[selectedOption.chain] || selectedOption.chain;
        const amount = parseFloat(selectedOption.amount).toFixed(Math.min(selectedOption.decimals, 8));
        
        detailsContainer.innerHTML = `
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
                <h4 style="margin: 0; color: #374151; font-size: 14px; font-weight: 500;">Send Payment:</h4>
                <button onclick="goBack()" style="color: #3b82f6; text-decoration: none; border: none; background: none; cursor: pointer; font-size: 14px;">← Back</button>
            </div>
            
            <div style="background: #f9fafb; border-radius: 8px; padding: 16px; margin-bottom: 16px;">
                <div style="display: flex; align-items: center; margin-bottom: 8px;">
                    <span class="chain-badge chain-${selectedOption.chain}" style="margin-right: 8px;">${chainName}</span>
                    <span style="font-weight: 500;">${selectedOption.symbol}</span>
                </div>
                <div style="font-size: 18px; font-weight: 600; color: #111827;">
                    ${amount} ${selectedOption.symbol}
                </div>
            </div>
            
            <div style="background: #f9fafb; border-radius: 8px; padding: 12px;">
                <div style="color: #6b7280; font-size: 12px; margin-bottom: 4px;">Send to address:</div>
                <div style="font-family: monospace; font-size: 14px; word-break: break-all; background: white; padding: 8px; border-radius: 4px; border: 1px solid #e5e7eb;">
                    ${selectedOption.address}
                </div>
            </div>
            
            <div style="text-align: center; margin-top: 16px;">
                <div style="color: #6b7280; font-size: 14px;">
                    🔍 Monitoring for payment...
                </div>
            </div>
        `;
    }

    window.goBack = function() {
        const optionsContainer = document.getElementById('payment-options');
        const detailsContainer = document.getElementById('payment-details');
        
        detailsContainer.style.display = 'none';
        optionsContainer.style.display = 'block';
        selectedOption = null;
    }

    function showSuccess() {
        container.innerHTML = `
            <div style="border: 1px solid #e5e7eb; border-radius: 8px; padding: 40px; background: white; text-align: center; max-width: 400px;">
                <div style="font-size: 48px; margin-bottom: 16px;">✅</div>
                <h3 style="margin: 0 0 8px 0; color: #111827; font-size: 18px; font-weight: 600;">Payment Successful!</h3>
                <p style="margin: 0; color: #6b7280;">Your payment has been confirmed.</p>
            </div>
        `;
        
        if (payment?.success_url) {
            setTimeout(() => {
                window.location.href = payment.success_url;
            }, 2000);
        }
    }

    function showExpired() {
        container.innerHTML = `
            <div style="border: 1px solid #e5e7eb; border-radius: 8px; padding: 40px; background: white; text-align: center; max-width: 400px;">
                <div style="font-size: 48px; margin-bottom: 16px;">⏰</div>
                <h3 style="margin: 0 0 8px 0; color: #111827; font-size: 18px; font-weight: 600;">Payment Expired</h3>
                <p style="margin: 0; color: #6b7280;">This payment link has expired.</p>
            </div>
        `;
    }

    function showError(message) {
        container.innerHTML = `
            <div style="border: 1px solid #e5e7eb; border-radius: 8px; padding: 40px; background: white; text-align: center; max-width: 400px;">
                <div style="color: #ef4444; font-size: 24px; margin-bottom: 16px;">⚠️</div>
                <h3 style="margin: 0 0 8px 0; color: #111827; font-size: 18px; font-weight: 600;">Error</h3>
                <p style="margin: 0; color: #6b7280;">${message}</p>
            </div>
        `;
    }

    async function checkStatus() {
        try {
            const response = await fetch(`${apiBaseUrl}/api/payments/${paymentId}/status`);
            const status = await response.json();
            
            if (status.status === 'paid') {
                showSuccess();
                stopStatusPolling();
            } else if (status.status === 'expired') {
                showExpired();
                stopStatusPolling();
            }
        } catch (error) {
            console.error('Error checking status:', error);
        }
    }

    let statusInterval;
    function startStatusPolling() {
        statusInterval = setInterval(checkStatus, 5000);
    }

    function stopStatusPolling() {
        if (statusInterval) {
            clearInterval(statusInterval);
            statusInterval = null;
        }
    }

    // Initialize
    createWidget();
    loadPayment();

    // Cleanup on page unload
    window.addEventListener('beforeunload', stopStatusPolling);
})();