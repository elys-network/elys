<!--
order: 4
-->

# Keeper

The `Keeper` module is a core component of the Elys Network's perpetual trading system. It encapsulates the logic for managing the state of perpetual contracts, handling transactions, and interfacing with other modules. This guide provides a detailed overview of the key functions and their roles within the module.

## Block Initialization

### Epoch and Interest Rate Initialization

The `BeginBlocker` function is triggered at the beginning of each block. It performs various checks and updates, including:

- **Epoch Check**: Determines if a new epoch has begun.
- **Interest Rate Computation**: Computes borrow interest rates for each pool.
- **Position Updates**: Updates the health of margin trading positions (MTPs) and processes funding fees.

### Processing Margin Trading Positions

`BeginBlockerProcessMTP` is responsible for updating individual MTPs, including:

- Calculating take profit liabilities and borrow rates.
- Handling borrow interest payments and funding fee collection.
- Checking MTP health and determining if force closure is necessary.

## Collateral and Position Calculations

### Minimum Collateral Calculation

`CalcMinCollateral` calculates the minimum collateral required to open a position based on leverage and the trading asset's price.

### Consolidating Collateral

`CalcMTPConsolidateCollateral` consolidates an MTP's collateral into the base currency, ensuring accurate position valuation.

### Borrow Interest Liabilities Calculation

`CalcMTPBorrowInterestLiabilities` computes the interest liabilities for an MTP, including unpaid borrow interest and new interest accrued within the current epoch.

### Take Profit Borrow Rate Calculation

`CalcMTPTakeProfitBorrowRate` calculates the borrow rate for an MTP's take profit custody, ensuring the rate meets minimum requirements.

### Take Profit Liabilities Calculation

`CalcMTPTakeProfitLiability` calculates the liabilities associated with an MTP's take profit custody in the base currency.

## Position Management Checks

### Maximum Open Positions Check

`CheckMaxOpenPositions` ensures the number of open MTPs does not exceed the maximum allowed limit.

### Pool Health Verification

`CheckPoolHealth` verifies the health of a pool, ensuring it is enabled and not closed, and that its health is above the minimum threshold.

### Same Asset Position Check

`CheckSameAssetPosition` checks if a user already has an open position with the same assets, allowing for consolidation.

### User Authorization Check

`CheckUserAuthorization` ensures a user is authorized to open a position, based on whitelisting if enabled.

## Position Closure

### Closing a Margin Trading Position

`Close` handles the closure of an MTP, repaying liabilities and returning any remaining collateral to the user.

### Closing a Long Position

`CloseLong` processes the closure of a long position, including handling borrow interest and estimating repayment amounts.

### Closing a Short Position

`CloseShort` processes the closure of a short position, similar to `CloseLong`, but adjusted for short positions.

## Position Opening

### Opening a New Margin Trading Position

`Open` processes the opening of a new MTP, determining the appropriate position type and validating assets.

### Consolidating Positions on Open

`OpenConsolidate` handles the consolidation of existing and new MTPs when the user already has an open position with the same assets.

### Opening a Long Position

`OpenLong` processes the opening of a new long position, including leverage calculation and collateral handling.

### Detailed Long Position Processing

`ProcessOpenLong` performs the detailed steps required to open a long position, including borrowing assets and updating pool health.

### Opening a Short Position

`OpenShort` processes the opening of a new short position, including leverage calculation and collateral handling.

### Detailed Short Position Processing

`ProcessOpenShort` performs the detailed steps required to open a short position, including borrowing assets and updating pool health.

## Position and Pool Handling

### Estimating and Repaying Liabilities

`EstimateAndRepay` estimates the repayment amount for a position and handles the repayment process.

### Estimating Swaps

`EstimateSwap` estimates the result of a token swap using the AMM's `CalcOutAmtGivenIn` function.

### Event Emissions

Defines various event emission functions used throughout the module, such as `EmitFundPayment`, `EmitForceClose`, and `EmitFundingFeePayment`.

### Forcibly Closing Positions

`ForceCloseLong` and `ForceCloseShort` handle the forced closure of positions when necessary, typically due to health or take profit conditions.

### Pool Management

`GetAmmPool`, `GetBestPool`, `SetPool`, and other functions handle the retrieval, initialization, and management of AMM pools.

### Borrow Interest Handling

`SettleBorrowInterest` and `SettleBorrowInterestPayment` manage the computation and payment of borrow interest for MTPs.

### Funding Fee Handling

`SettleFundingFeeCollection` and `HandleFundingFeeDistribution` manage the collection and distribution of funding fees.

### Health Updates

`GetMTPHealth` and `UpdatePoolHealth` ensure the health of MTPs and pools are calculated and maintained.

### Utility Functions

Various utility functions such as `GetMarketAssetPrice`, `PositionEstimation`, and `SwapAssets` provide necessary support for position and pool management.

## Governance and Whitelisting

### Message Server Implementation

Defines the `MsgServer` struct and its implementation, handling various message types such as opening and closing positions, whitelisting, and parameter updates.

### Governance Proposals

`UpdateParams` handles updating the module's parameters through a governance proposal.

### Whitelisting Management

`Whitelist` and `Dewhitelist` manage the addition and removal of addresses from the whitelist, ensuring only authorized users can open positions.

## Queries

### Querying Positions

`GetPositions`, `GetPositionsByPool`, and `GetPositionsForAddress` provide endpoints to retrieve open positions.

### Querying Pool and Module Status

`GetStatus`, `Pools`, `Pool`, and `Params` provide endpoints to retrieve the status of the module, including pool details and module parameters.

### Querying Whitelist

`GetWhitelist` and `IsWhitelisted` provide endpoints to retrieve and check whitelisted addresses.

### Position Opening Estimation

`OpenEstimation` provides an estimation for opening a position, including collateral requirements, estimated PnL, and liquidation prices.
