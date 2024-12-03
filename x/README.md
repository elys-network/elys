# Elys Network Modules

This document provides a summary of each module available under the `x` folder of the Elys Network, along with links to the respective module specification documents. Each module is designed to enhance the functionalities of the Elys Network by adding specific capabilities like automated market making, staking, liquidity provisioning, asset management, and more. For detailed information, click on the links provided for each module.

## Modules Overview

### 1. [AccountedPool Module](accountedpool)

The `accountedpool` module manages and maintains the state of accounted pools within the Elys Network. It ensures accurate accounting of pool balances, integrating with other modules like AMM and Perpetual pools.

- **Features**: Accurate pool management, integration with AMM and Perpetual pools, query services.
- **[Module Spec](accountedpool/spec/README.md)**

### 2. [AMM Module](amm)

The Automated Market Maker (AMM) module supports liquidity pools designed for assets with significant price variation and those with stable prices through AMM and Oracle pools.

- **Features**: Liquidity provision, oracle integration, risk management.
- **[Module Spec](amm/spec/README.md)**

### 3. [Asset Profile Module](assetprofile)

The `assetprofile` module manages asset properties, defining parameters and handling Inter-Blockchain Communication (IBC) integration.

- **Features**: Asset management, IBC integration, parameter management.
- **[Module Spec](assetprofile/spec/README.md)**

### 4. [Burner Module](burner)

The `burner` module allows for automatic burning of native tokens at regular intervals, depending on the Epochs module.

- **Features**: Token burning, integration with the Epochs module.
- **[Module Spec](burner/spec/README.md)**

### 5. [Commitment Module](commitment)

The `commitment` module manages token commitments, including staking, vesting, and locking of tokens.

- **Features**: Token commitment, staking, vesting schedules, dynamic parameter updates.
- **[Module Spec](commitment/spec/README.md)**

### 6. [Epochs Module](epochs)

The `epochs` module provides a generalized epoch interface, allowing other modules to execute tasks at specified time intervals.

- **Features**: Time-based task execution, generalized epoch signaling.
- **[Module Spec](epochs/spec/README.md)**

### 7. [eStaking Module](estaking)

The `estaking` module extends basic staking functionalities by adding advanced reward management, staking parameter updates, and Eden token mechanics.

- **Features**: Advanced reward distribution, Eden token management, staking parameter updates.
- **[Module Spec](estaking/spec/README.md)**

### 8. [LeverageLP Module](leveragelp)

The `leveragelp` module allows users to add liquidity in leverage in AMM pools to enhance their rewards while ensuring safety through health checks and liquidation mechanisms.

- **Features**: Leveraged liquidity, health checks, dynamic parameter updates.
- **[Module Spec](leveragelp/spec/README.md)**

### 9. [Masterchef Module](masterchef)

The `masterchef` module manages liquidity provider rewards, external incentives, and dynamically updates staking parameters.

- **Features**: Reward management, external incentives, dynamic parameter updates.
- **[Module Spec](masterchef/spec/README.md)**

### 10. [Oracle Module](oracle)

The `oracle` module provides decentralized price feeds and manages asset information by utilizing multiple sources to ensure reliability.

- **Features**: Decentralized price feeds, asset information management, price feeder control.
- **[Module Spec](oracle/spec/README.md)**

### 11. [Parameter Module](parameter)

The `parameter` module manages and maintains key configuration parameters within the Elys Network, allowing for dynamic and controlled adjustments.

- **Features**: Dynamic parameter management, query services, controlled updates.
- **[Module Spec](parameter/spec/README.md)**

### 12. [Perpetual Module](perpetual)

The `perpetual` module facilitates perpetual trading, allowing users to open and close leveraged positions without expiry dates, with various safety and health checks.

- **Features**: Perpetual trading, leverage management, safety factor, liquidation mechanisms.
- **[Module Spec](perpetual/spec/README.md)**

### 13. [Stablestake Module](stablestake)

The `stablestake` module manages stable staking functionalities, including borrowing and lending mechanisms, interest rate management, and debt handling.

- **Features**: Borrowing and lending management, interest rate updates, debt management.
- **[Module Spec](stablestake/spec/README.md)**

### 14. [Tier Module](tier)

The `tier` module manages the tier membership system, providing loyal users with discounts and benefits across all services available in the Elys Network.

- **Features**: Tier membership management, user discounts, service-wide benefits.
- **[Module Spec](tier/spec/README.md)**

### 15. [Tokenomics Module](tokenomics)

The `tokenomics` module manages the economic and incentive mechanisms of the network, including airdrops and inflation.

- **Features**: Airdrop management, inflation handling, dynamic parameter updates.
- **[Module Spec](tokenomics/spec/README.md)**

### 16. [TradeShield Module](tradeshield)

The `tradeshield` module provides functionalities for creating and managing various types of market orders, including spot and perpetual orders, and handles order execution through off-chain agents to optimize performance.

- **Features**: Order creation and cancellation, off-chain execution agents, penalty and reward systems.
- **[Module Spec](tradeshield/spec/README.md)**

### 17. [Transferhook Module](transferhook)

The `transferhook` module provides advanced functionality for handling IBC transfers, integrating AMM interactions for efficient token transfers.

- **Features**: AMM integration, parameter management, query services.
- **[Module Spec](transferhook/spec/README.md)**

## References

- [Elys Network Documentation](https://docs.elys.network)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [GitHub Repository for Elys Network](https://github.com/elys-network/elys)

For detailed information about each module, refer to their respective documentation linked above.