# Off-Chain Logic for Position Health and Stop-Loss Checks

## Check Position Health Logic

### Requirements:

1. **Position and Debt Objects:**
   - `Position` object contains `lpAmount` and other relevant details.
   - `Debt` object contains information about the borrowed amount, interest, etc.

2. **Position Health Calculation:**
   - **Formula:** Position Health = (Value of `lpAmount` in USDC) / (Liability in USDC).
   - Ref: [Position Health Calculation](https://github.com/elys-network/elys/blob/90d8b7e326dd424a622bb7977e5f2bf1f2a0f1ad/x/leveragelp/keeper/position.go#L388).
   - **Value of `lpAmount` in USDC:**
     - Calculated using the AMM (Automated Market Maker) pool.
   - **Liability Calculation:**
     - Include interest stacked over time on the debt.
     - Use `Debt` object and `stablestake` parameters to calculate interest stacked.
     - Ref: [Debt Interest Calculation](https://github.com/elys-network/elys/blob/90d8b7e326dd424a622bb7977e5f2bf1f2a0f1ad/x/stablestake/keeper/debt.go#L63).
     - After calculating use **Formula:** Liability = Debt.Borrowed + Debt.InterestStacked - Debt.InterestPaid.

     - If position health less than safety factor, trigger `closePositions` function.
 
## Check Stop-Loss Conditions

### Requirements:

1. **Stop-Loss Trigger:**
   - Trigger the `closePosition` function if the `lpAmount` price reaches the stop-loss price.
   - **Required Data:**
     - `lpAmount` price (obtained using the AMM module).
     - Stop-loss price (present in the `Position` object).

### Summary:

1. Calculate the current value of `lpAmount` in USDC using the AMM pool.
2. Calculate the liability by accounting for borrowed amount and interest (stacked and paid).
3. Calculate the Position Health: Position Health = (Value of `lpAmount` in USDC) / (Liability in USDC).
4. If position health less than safety factor, trigger `closePositions` function.
5. Check if the current `lpAmount` price has reached the stop-loss price, and if so, trigger the `closePositions` function.

### closePositions function
- `closePositions(liquidate: [](address, u64), stopLoss: [](address, u64))`
  - liquidate: list of position ids that needs to be liquidated, function will verify and liquidate. Tuple needs to have user address and position id.
  - stopLoss: list of position ids that needs to be closed, function will check for stop loss and close position. Tuple needs to have user address and position id.
