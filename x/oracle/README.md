# Oracle Module

The `oracle` module provides decentralized price feeds and asset information through an integrated price feeder functionality within the ELYS binary. The module handles price updates at the ABCI level and allows node operators to participate in feeding price data by enabling the price feeder configuration. This integration enhances the efficiency and reliability of price-related operations across the network by tightly coupling the price feeding process with the blockchain's consensus mechanism.

## Important Notice

The Oracle module implementation has been moved to the [Ojo repository](https://github.com/ojo-network/ojo/tree/main/x/oracle). This document provides an overview of how the Oracle module is used within the ELYS network.

## Integration with ELYS

The price feeder functionality, which was previously a separate off-chain agent, is now fully integrated into the ELYS binary. This means that node operators can enable price feeding capabilities directly from their node by using the following configuration:

```bash
# Enable price feeder when starting the node
elysd start --price-feeder.enable=true --price-feeder.config=/path/to/price_feeder_config.toml
```

This integration allows any node operator to participate in feeding price and external liquidity data to the Oracle module.

## Price Updates

Price updates happen at the ABCI (Application Blockchain Interface) level rather than through external messages. This means the price feeding process is more tightly integrated with the blockchain's consensus mechanism.

## Available Queries

The following queries are available for retrieving price information:

- `show-price`: Get the price for a specific asset
- `list-price`: Get a list of all prices
- `exchange-rate`: Get the current exchange rate for a specific asset
- `exchange-rates`: Get a summary of exchange rates for all assets

## Implementation Details

For detailed implementation specifications of the Oracle module, including:
- Voting procedures
- Reward mechanisms
- Slashing conditions
- State management
- Message handling
- Events

Please refer to the [Ojo Oracle Module Documentation](https://github.com/ojo-network/ojo/tree/main/x/oracle).