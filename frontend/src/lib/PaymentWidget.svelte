<script>
	import { onMount, onDestroy } from 'svelte';
	import PaymentOption from './PaymentOption.svelte';
	import QRCode from './QRCode.svelte';

	export let paymentId;

	let payment = null;
	let selectedOption = null;
	let loading = true;
	let error = null;
	let statusInterval;

	const API_BASE_URL = window.API_BASE_URL || 'http://localhost:8080';

	async function fetchPayment() {
		try {
			const response = await fetch(`${API_BASE_URL}/api/payments/${paymentId}`);
			if (!response.ok) throw new Error('Payment not found');
			
			payment = await response.json();
			loading = false;

			if (payment.status === 'paid') {
				redirectToSuccess();
				return;
			}

			if (payment.status === 'expired') {
				error = 'Payment has expired';
				return;
			}

			// Start checking payment status
			startStatusPolling();
		} catch (err) {
			error = err.message;
			loading = false;
		}
	}

	async function checkPaymentStatus() {
		try {
			const response = await fetch(`${API_BASE_URL}/api/payments/${paymentId}/status`);
			if (!response.ok) return;
			
			const status = await response.json();
			if (status.status === 'paid') {
				redirectToSuccess();
			} else if (status.status === 'expired') {
				error = 'Payment has expired';
				stopStatusPolling();
			}
		} catch (err) {
			console.error('Error checking payment status:', err);
		}
	}

	function startStatusPolling() {
		statusInterval = setInterval(checkPaymentStatus, 5000);
	}

	function stopStatusPolling() {
		if (statusInterval) {
			clearInterval(statusInterval);
			statusInterval = null;
		}
	}

	function redirectToSuccess() {
		stopStatusPolling();
		if (payment?.success_url) {
			window.location.href = payment.success_url;
		} else {
			// Show success message
			error = null;
			payment = { ...payment, status: 'paid' };
		}
	}

	function selectOption(option) {
		selectedOption = option;
	}

	function getChainName(chain) {
		switch (chain) {
			case 'ethereum': return 'Ethereum';
			case 'solana': return 'Solana';
			case 'ton': return 'TON';
			default: return chain;
		}
	}

	function formatAmount(amount, decimals) {
		return parseFloat(amount).toFixed(Math.min(decimals, 8));
	}

	onMount(() => {
		fetchPayment();
	});

	onDestroy(() => {
		stopStatusPolling();
	});
</script>

<div class="bg-white rounded-lg shadow-lg p-6 w-full max-w-md">
	{#if loading}
		<div class="text-center py-8">
			<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500 mx-auto mb-4"></div>
			<p class="text-gray-600">Loading payment...</p>
		</div>
	{:else if error}
		<div class="text-center py-8">
			<div class="text-red-500 text-xl mb-2">⚠️</div>
			<h3 class="text-lg font-medium text-gray-900 mb-2">Payment Error</h3>
			<p class="text-gray-600">{error}</p>
		</div>
	{:else if payment?.status === 'paid'}
		<div class="text-center py-8">
			<div class="text-green-500 text-4xl mb-4">✅</div>
			<h3 class="text-lg font-medium text-gray-900 mb-2">Payment Successful!</h3>
			<p class="text-gray-600">Your payment has been confirmed.</p>
		</div>
	{:else if payment}
		<div>
			<h2 class="text-xl font-semibold text-gray-900 mb-4">Complete Your Payment</h2>
			
			<div class="mb-6">
				<div class="bg-gray-50 rounded-lg p-4">
					<div class="flex justify-between items-center">
						<span class="text-sm text-gray-600">Amount:</span>
						<span class="text-lg font-semibold">${payment.amount} {payment.currency}</span>
					</div>
				</div>
			</div>

			{#if !selectedOption}
				<div class="mb-6">
					<h3 class="text-sm font-medium text-gray-900 mb-3">Choose Payment Method:</h3>
					<div class="space-y-2">
						{#each payment.options as option (option.id)}
							<PaymentOption {option} onClick={() => selectOption(option)} />
						{/each}
					</div>
				</div>
			{:else}
				<div class="mb-6">
					<div class="flex items-center justify-between mb-4">
						<h3 class="text-sm font-medium text-gray-900">Send Payment:</h3>
						<button 
							class="text-sm text-primary-600 hover:text-primary-700"
							on:click={() => selectedOption = null}
						>
							← Back
						</button>
					</div>
					
					<div class="bg-gray-50 rounded-lg p-4 mb-4">
						<div class="flex items-center mb-2">
							<span class="chain-badge chain-{selectedOption.chain} mr-2">
								{getChainName(selectedOption.chain)}
							</span>
							<span class="font-medium">{selectedOption.symbol}</span>
						</div>
						<div class="text-lg font-semibold text-gray-900">
							{formatAmount(selectedOption.amount, selectedOption.decimals)} {selectedOption.symbol}
						</div>
					</div>

					<div class="text-center mb-4">
						<QRCode value={selectedOption.address} size={200} />
					</div>

					<div class="bg-gray-50 rounded-lg p-3">
						<div class="text-xs text-gray-600 mb-1">Send to address:</div>
						<div class="font-mono text-sm break-all bg-white p-2 rounded border">
							{selectedOption.address}
						</div>
					</div>

					<div class="mt-4 text-center">
						<div class="animate-pulse text-sm text-gray-600">
							🔍 Monitoring for payment...
						</div>
					</div>
				</div>
			{/if}
		</div>
	{/if}
</div>