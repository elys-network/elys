# Off-Chain Logic for Position Health and Stop-Loss Checks

## Check Position Health Logic

### Requirements:

1. **Position Objects:**
    - `Position` object contains `MtpHealth`.
    - If `MtpHealth` is less than safety factor, trigger `closePositions` function.
    - Safety factor can be retrieved from module parameters

## Check Stop-Loss Conditions

### Requirements:

1. **Stop-Loss Trigger:**
  - Trigger the `closePosition` function if the `TradingAsset` price reaches the stop-loss price.
  - **Required Data:**
    - `TradingAsset` price (obtained using the oracle module).
    - Stop-loss price (present in the `Position` object).
  - For long position: If `tradingAssetPrice <= mtp.StopLossPrice` then trigger closePosition message.
  - For short position: If `tradingAssetPrice >= mtp.StopLossPrice` then trigger closePosition message.

## Check Take-Profit Conditions

### Requirements:

1. **Take-Profit Trigger:**
  - Trigger the `closePosition` function if the `tradingAsset` price reaches the take-profit price.
  - **Required Data:**
    - `TradingAsset` price (obtained using the oracle module).
    - `TakeProfitPrice` price (present in the `Position` object).
  - For long position: If `tradingAssetPrice >= mtp.TakeProfitPrice` then trigger closePosition message.
  - For short position: If `tradingAssetPrice <= mtp.TakeProfitPrice` then trigger closePosition message.

  - Note: only oracle pools are used by Perpetual, therefore the bot should always relying on the Oracle module for retrieving asset price.

### closePositions function
- `closePositions(liquidate: [](address, u64), stopLoss: [](address, u64))`
  - liquidate: list of position ids that needs to be liquidated, function will verify and liquidate. Tuple needs to have user address and position id.
  - stopLoss: list of position ids that needs to be closed, function will check for stop loss and close position. Tuple needs to have user address and position id.
