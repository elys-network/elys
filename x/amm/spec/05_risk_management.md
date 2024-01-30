<!--
order: 5
-->

# Risk management

## Risk management scenarios & solutions on hybrid AMM model

- Arbitrage opportunity?
- Permanent loss for LPs for providing better swap?
- Utilization rate vs loss?
- Oracle accuracy assumption - TODO: calculate on loss amount when oracle has 1% price difference with Osmosis.
- How much do LPs loss when price go down and recover to original?
- Need to have a way to stop deposit/withdrawal of specific token within the vault?
- Need to have a way to stop overall deposit/withdrawal on the vault?
- When oracle has large difference between several platforms, it's not trustworthy, and `oracle` based pool will need ask high fees so that it does not give more than the minimum price or just redirect the operation to normal Elys AMM.

## To consider

- Protect LP from volatile markets and market arbs
- Oracle difference % from sudden dump on third party market
- Weights are balanced by the incentive structure to hit close to the target weights
- Unit test to check final slippage on example values
- Low liquidity coin and high liquidity coin (low liquidity coin won't be able to use full oracle price)
- Once the user buys the major molecule TKN, it is automatically “committed” in the commitment molecule to receive rewards
- LPs should not lose underlying value as long as underlying tokens are recovered back to original status

## Automatic pool parameters adjustment from asset volatility

```protobuf
message AssetVolatility {
    string asset = 1;
    string volatility = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
}
```

```go
Major molecule tokens = (USD value of deposit Asset - fees)/Price of Major Molecule.

Initial Major molecule token price = $1

Number of The fee is always swapped to USDC and sent to the major molecule fee wallet which stores all the revenue.

Price of Major molecule token at any point = (USD value of all assets within the major molecule +/- Perpetual Gains/Losses)/circulating supply of Major Molecule Tokens
```

### Slippage

Slippage is applied to secure liquidity providers' bonded liquidity.

There should be well-formed slippage calculator considering oracle. (Consider GMX model)

We will introduce `external_liquidity` parameter to provide as less as slippage considering external exchanges.

One option is to introduce dynamically configurable slippage by governance of molecule token.

Volatility could be configured for individual assets, and slippage could be configured for pool itself.

For 4 asset pool (ETH, BTC, ATOM, USDC), different slippage could be executed based on following option.

Volatility of swapping assets + pool slippage + target weight change

TODO: should have exact Maths formula for calculating slippage.

## Asset weight recovery once it's broken

Following items are applied to keep the weight not broken

1. Slippage on swap when imbalanced
2. More fees on deposit when imbalanced

But this is not enough incentivization to recover the broken weight.

Therefore, we additionally introduce weight recovery treasury

1. 10% of fees are put in weight recovery treasury
2. Once weight is broken, users who recover the weight by doing swap, or adding new lp are incentivized.

TODO: determine the distribution mechanism of weight recovery treasury.

One way of recovering the imbalanced weight is to use $10K fund from the team to recover imbalanced $100K.
E.g. JUNO has $100K more and it's imbalanced, team could swap $10K worth of USD to JUNO on Elys, swap received JUNO to USD on Osmosis and repeated the process 10 times to recover whole weight.
During the execution, fees will be spent, but weight should be recovered spending no more than $1K.
