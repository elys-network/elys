# Elys Leverage LP Module

## Contents

1. **[Concepts](01_concepts.md)**
2. **[Usage](02_usage.md)**
3. **[Keeper](03_keeper.md)**
4. **[Protobuf Definitions](04_protobuf_definitions.md)**
5. **[Functions](05_functions.md)**

## References

Resources:

- [Elys Network Documentation](https://docs.elys.network)
- [Cosmos SDK Documentation](https://docs.cosmos.network)
- [GitHub Repository for Elys Network](https://github.com/elys-network/elys)

## Overview

The `leveragelp` module in the Elys Network enables users to provide leveraged liquidity in AMM pools. By utilizing collateral and leverage, users can enhance their returns from liquidity

provisioning while the module ensures the safety and stability of these leveraged positions through rigorous health checks and liquidation mechanisms.

## Key Features

- **Leveraged Liquidity Provision**: Allow users to open leveraged positions in AMM pools.
- **Health Checks and Liquidation**: Regularly check the health of positions and liquidate unhealthy ones to maintain stability.
- **Whitelist Management**: Manage access through address whitelisting.
- **Flexible Parameter Management**: Update leverage limits, pool thresholds, and other parameters dynamically.

For more detailed information, please refer to the individual sections listed in the contents above.
