<!--
order: 5
-->

# Endpoints

## Gov Proposals

- CreateVault(assets, ratios)
- UpdateVaultConfig(vaultId, targetWeights) - Pool params will need to be governed by governance.

## Msg endpoints

- Deposit(assets) - calculate slippage based on weight change
- Swap(asset->target_asset) - calculate slippage based on weight change
- Withdraw(lp->target_asset) - calculate slippage based on weight change

## Query endpoints

- QueryVaults - Vault configs
- QueryVault(vaultId) - LP price, vault config, vault TVL, fees collected
- EstimatedSwapOutAmount(SwapCoin, outToken) -> amount
- EstimatedLPTokenAmountAfterDeposit(depositCoins) -> amount
- EstimatedWithdrawAmount(lpTokenAmount, outToken) -> outTokenAmount

## Epoch actions

- Fee distribution
- Trading volume track?
- Liquidity P/L track?
