<script>
	export let option;
	export let onClick;

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

	function getTokenIcon(symbol) {
		switch (symbol) {
			case 'ETH': return '🔷';
			case 'SOL': return '🟣';
			case 'TON': return '💎';
			case 'USDC': return '💵';
			case 'USDT': return '💵';
			default: return '💰';
		}
	}
</script>

<button 
	class="payment-option w-full text-left" 
	on:click={onClick}
>
	<div class="flex items-center justify-between">
		<div class="flex items-center">
			<span class="text-xl mr-3">{getTokenIcon(option.symbol)}</span>
			<div>
				<div class="flex items-center">
					<span class="font-medium text-gray-900 mr-2">{option.symbol}</span>
					<span class="chain-badge chain-{option.chain}">
						{getChainName(option.chain)}
					</span>
				</div>
				<div class="text-sm text-gray-600">
					{formatAmount(option.amount, option.decimals)} {option.symbol}
				</div>
			</div>
		</div>
		<div class="text-gray-400">
			→
		</div>
	</div>
</button>