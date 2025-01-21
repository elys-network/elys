# TradeShield Module

The purpose of this document is to specify the design and implementation of the new TradeShield native module for the Elys Network blockchain. The goal is to migrate the existing TradeShield (TS) CosmWasm contract functionalities to a native module to reduce the complexity and computational overhead, contributing to the reduction of block times to 1-2 seconds.

## Objectives

- **Reduce ABCI Logic:** Migrate TradeShield CosmWasm contract to a native Elys module to reduce ABCI logic, particularly in the `x/clock` end blocker, and move closer to achieving 1-2 seconds block time.
- **Simplify Codebase:** Eliminate bindings logic from the Elys repository related to CW contracts, simplifying the overall codebase.
- **Optimize Performance:** Offload computationally heavy order execution checks from the chain to off-chain agents (bots), reducing block processing time.

## Module Parameters

The TradeShield module will have the following parameters:

1. **Reward Percentage:** The percentage of the operation value that will be given to the participant as a reward for successfully executed orders. (Example: 0.1%)
2. **Margin Error Rate:** The margin of error for order execution. If the order price is within this margin above the market price, the order will not be executed. (Example: 1%)
3. **Minimum Deposit:** The minimum amount of Elys tokens that must be deposited along with the transaction. This deposit will be forfeited if any order in the submitted list does not meet the execution criteria.

## Functionalities

The TradeShield module will replicate and extend the functionalities currently available in the TradeShield contract. The following functionalities will be supported:

### Order Creation and Fund Management

When a new order is created, the order amount is automatically moved from the user's address to a new order-specific address. This locks the funds and prevents them from being used in other transactions while the order is active. Users can retrieve their locked funds by canceling the order, but only if the order is still pending execution. Once an order is executed, the locked funds are used either for spot order AMM swaps or to open perpetual positions.

1. **Spot Market Orders**

   - Create market order (buy, sell)

2. **Spot Limit Orders**

   - Create spot limit order (buy, sell)
   - Create stop-loss order

3. **Perpetual Market Orders**

   - Create perpetual market open order

4. **Perpetual Limit Orders**

   - Create perpetual limit open order
   - Create perpetual stop-loss order
   - Create perpetual limit close order

### Order Cancellation

1. **Spot Orders**

   - Cancel spot limit order

2. **Perpetual Orders**
   - Cancel perpetual limit order

### Querying

1. **Orders**

   - Query orders by various parameters (e.g., market, type, status)

2. **Order by ID**
   - Query specific order by ID

## Execution Flow

### Off-Chain Execution Agents (Bots)

- **Order Execution:** Off-chain agents will be responsible for submitting a list of order IDs that need to be executed. This eliminates the need for the blockchain to traverse all orders at every block, significantly reducing computational load.
- **Chain Verification:** Despite offloading execution logic, the blockchain will still verify that the execution criteria of submitted orders are met before executing them.

### Permissionless Execution

- **Open Participation:** Any address can participate in submitting orders for execution. The process is designed to be permissionless, similar to the liquidation system.

### Incentive and Penalty

- **Penalties for Invalid Submissions:**

  - Participants must deposit a minimum amount of Elys tokens along with their transaction.
  - If any order in the submitted list does not meet the execution criteria, the deposit is forfeited.
  - If all orders meet the criteria, the deposit is refunded.

- **Rewards for Valid Submissions:**
  - Participants will receive a percentage of the operation value as a reward for successfully executed orders.
  - Example: For a spot limit buy order, the order is executed 1% before the order price, and the 1% value is given to the participant as a reward, with the remaining output sent to the order creator.

### Margin of Error

To account for brief fluctuations in asset price, a margin of error is introduced:

- **Margin Error Rate**: Set at 1%.
- **Application**:
  - **Order Price**: If the order price is within 1% above the market price, the order will not be executed, and the deposit will be returned without penalty.

## Integration with Indexer

- **Indexer Utilization:**
  - The module will integrate with an indexer to reduce the number of queries made against the network nodes.
  - This integration is intended to optimize performance by allowing efficient retrieval of order data without overloading the network nodes.
