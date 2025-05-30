<script>
	import PaymentWidget from '$lib/PaymentWidget.svelte';
	import { onMount } from 'svelte';

	let paymentId = '';

	onMount(() => {
		// Get payment ID from URL params or global variable
		const urlParams = new URLSearchParams(window.location.search);
		paymentId = urlParams.get('payment_id') || window.PAYMENT_ID || '';
	});
</script>

<svelte:head>
	<title>Crypto Payment</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 flex items-center justify-center p-4">
	{#if paymentId}
		<PaymentWidget {paymentId} />
	{:else}
		<div class="text-center">
			<h1 class="text-2xl font-bold text-gray-900 mb-4">Payment Widget</h1>
			<p class="text-gray-600">No payment ID provided</p>
		</div>
	{/if}
</div>