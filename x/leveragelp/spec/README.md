# leveragelp module

`leveragelp` module provides interface for users to add liquidity in leverage in AMM pools to get more rewards. By utilizing collateral and leverage, users can enhance their returns from liquidity provisioning while the module ensures the safety and stability of these leveraged positions through rigorous health checks and liquidation mechanisms.

The underlying mechanism is when a user would like to open a leveraged position, module borrows `USDC` from stablestake module and put it as liquidity along with collateral amount received from the user.

To make it secure, the module monitors position value and if it goes down below threshold, it is force closing the position and returning all the borrowed `USDC` back to `stablestake` module along with stacked interest.

## Key Features

- **Leveraged Liquidity Provision**: Allow users to open leveraged positions in AMM pools.
- **Health Checks and Liquidation**: Regularly check the health of positions and liquidate unhealthy ones to maintain stability.
- **Whitelist Management**: Manage access through address whitelisting.
- **Flexible Parameter Management**: Update leverage limits, pool thresholds, and other parameters dynamically.

For more detailed information, please refer to the individual sections listed in the contents below.

## Contents

1. **[Concepts](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Keeper](03_keeper.md)**
4. **[Endpoints](04_endpoints.md)**
5. **[CLI](05_cli.md)**
6. **[Position Workflow](06_position_workflow.md)**
7. **[Pool Status](07_pool_status.md)**
