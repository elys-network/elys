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
   - Trigger the `closePosition` function if the `custody` price reaches the stop-loss price.
   - **Required Data:**
     - `custody` price (obtained using the AMM module).
     - Stop-loss price (present in the `Position` object).

### closePositions function
- `closePositions(liquidate: [](address, u64), stopLoss: [](address, u64))`
  - liquidate: list of position ids that needs to be liquidated, function will verify and liquidate. Tuple needs to have user address and position id.
  - stopLoss: list of position ids that needs to be closed, function will check for stop loss and close position. Tuple needs to have user address and position id.
