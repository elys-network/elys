# Perpetual module

The `perpetual` module in the Elys Network facilitates perpetual trading, allowing users to open and close leveraged positions on various assets without expiry dates. Positions, either long or short, are determined by the leverage applied and the amount of collateral provided. Leverage allows traders to borrow funds to open larger positions, amplifying potential profits and losses.

Collateral acts as a security deposit to cover potential losses, and if it falls below the maintenance margin, the position may be liquidated. The safety factor is a threshold to keep positions open, and position health indicates the risk level of a position. Liquidation occurs when the position's value drops below the safety factor, ensuring platform stability.

The funding rate is a periodic payment exchanged between long and short positions to align contract prices with underlying asset prices. Whitelisting controls access to the module, allowing only trusted participants to trade. The module has configurable parameters, such as leverage limits and interest rates, adjustable through governance proposals to maintain efficiency and security.

## Key Features

- **Perpetual Trading**: Trade assets with leverage without expiry dates.
- **Positions**: Open long or short leveraged positions based on collateral and leverage.
- **Leverage**: Borrow funds to open larger positions, amplifying profits and losses.
- **Collateral**: Security deposit to cover potential losses and prevent liquidation.
- **Safety Factor**: Threshold to keep positions open and prevent liquidation.
- **Position Health**: Metric indicating the risk level of a position.
- **Liquidation**: Forcibly close positions below the safety factor to prevent further losses.

For more detailed information, please refer to the individual sections listed in the contents below.

## Contents

1. **[Concepts](01_concepts.md)**
2. **[Mechanisms](02_mechanism.md)**
3. **[Usage](03_usage.md)**
4. **[Keeper](04_keeper.md)**
5. **[Protobuf Definitions](05_protobuf_definitions.md)**
6. **[Functions](06_functions.md)**
7. **[Pool Status](07_pool_status.md)**
