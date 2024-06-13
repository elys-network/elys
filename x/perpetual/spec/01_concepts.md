<!--
order: 1
-->

# Concepts

The `perpetual` module in the Elys Network facilitates perpetual trading, allowing users to open and close positions on various assets. This guide provides an overview of the key concepts and mechanisms involved in the perpetual module.

## Key Concepts

### Perpetual Trading

Perpetual trading allows users to trade assets with leverage without the need for expiry dates on their positions. This type of trading is popular in the cryptocurrency market as it offers flexibility and continuous trading opportunities.

### Positions

A position in perpetual trading represents an open contract on an asset. Positions can be either long (betting that the asset price will rise) or short (betting that the asset price will fall). The size of the position is determined by the leverage applied and the amount of collateral provided.

### Leverage

Leverage allows traders to open positions larger than their collateral by borrowing funds. For example, with 5x leverage, a trader can open a position five times the size of their collateral. While leverage amplifies potential profits, it also increases potential losses.

### Collateral

Collateral is the amount of assets locked to open a leveraged position. It acts as a security deposit to cover potential losses. If the value of the position falls below the maintenance margin, the position may be liquidated to prevent further losses.

### Safety Factor

The safety factor is the threshold to keep a position open. If the position health value falls below this threshold, the position may be subject to liquidation.

### Position Health

Position health is a metric that indicates the risk level of a position. It is calculated based on the position's value, leverage, and collateral. Monitoring position health helps prevent excessive risk-taking and ensures the stability of the perpetual trading system.

### Liquidation

Liquidation occurs when a position's value drops below the safety factor. The position is forcibly closed to prevent further losses. This mechanism ensures the stability of the perpetual trading system and protects the platform from significant losses.

### Funding Rate

The funding rate is a periodic payment exchanged between long and short positions. It ensures that the perpetual contract price closely tracks the underlying asset price. When the funding rate is positive, longs pay shorts, and when it is negative, shorts pay longs.

### Whitelisting

Whitelisting is a mechanism to control access to the perpetual module. Only addresses that are whitelisted can participate in trading. This helps maintain the security and integrity of the platform by allowing only trusted participants.

### Parameters

The perpetual module has several configurable parameters that govern its operation. These include leverage limits, interest rates, maintenance margins, and more. These parameters can be updated through governance proposals to adapt to changing market conditions and ensure the module's efficiency and security.

## Transaction Commands

The perpetual module supports various transaction commands for managing positions and parameters. Key commands include:

- `open`: Opens a new perpetual position.
- `close`: Closes an existing perpetual position.
- `whitelist`: Adds an address to the whitelist.
- `dewhitelist`: Removes an address from the whitelist.
- `update-params`: Updates the module parameters through a governance proposal.

## Query Commands

Users can query various aspects of the perpetual module using the following commands:

- `params`: Retrieves the current parameters of the module.
- `get-positions`: Lists all open positions.
- `get-positions-by-pool`: Lists positions for a specific pool.
- `get-positions-for-address`: Lists positions for a specific address.
- `get-status`: Retrieves the current status of the module.
- `get-whitelist`: Lists all whitelisted addresses.
- `is-whitelisted`: Checks if a specific address is whitelisted.
- `list-pool`: Lists all available pools.
- `show-pool`: Retrieves details of a specific pool.
- `get-mtp`: Retrieves details of a specific margin trading position (MTP).
- `open-estimation`: Provides an estimation for opening a position.
