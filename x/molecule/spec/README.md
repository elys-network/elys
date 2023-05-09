# Molecule module

Resources:
https://www.notion.so/5-Molecule-Token-Adding-Liquidity-to-a-Multi-Asset-Pool-1fc49e859f96457a9469d9c9ed5d06cf

Vault contract (Main logic seems to be here)
https://arbiscan.io/address/0x489ee077994B6658eAfA855C308275EAd8097C4A#code

Reward Router - https://arbiscan.io/address/0xB95DB5B167D75e6d04227CfFFA61069348d271F5#code

GLP manager - https://arbiscan.io/address/0x3963FfC9dff443c2A94f21b129D429891E32ec18#code

GLP documenet

- https://gmxio.gitbook.io/gmx/glp
- https://app.gmx.io/#/buy_glp

## State

### Vault

```protobuf
message AssetVolatility {
    string asset = 1;
    string volatility = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
}
message AssetWeight {
    string asset = 1;
    string weight = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ];
}
message Vault {
    repeated AssetWeight weights = 1; // default: 22.5% ETH, 20% BTC, 22.5% ATOM, 35% USDC
    string slippage_discount = 2 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ]; // slippage discount is setting the slippage discount where slippage calculation is done based on elys liquidity
    string default_swap_fee = 3 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ]; // fee amount on swap without considering weight - 0.01%
    string slippage_fee = 4 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ]; // fee amount for weight change after the operation - fees will vary between 0.01% to 1%
    string lp_fee_portion = 5 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ]; // fees that is sent to liquidity providers - 65%
    string stake_fee_portion = 6 [
        (gogoproto.customtype) = "github.com/cosmos/cosmos-sdk/types.Dec",
        (gogoproto.nullable) = false
    ]; // fees that is sent to stakers - 35%
}
```

```go
Major molecule tokens = (USD value of deposit Asset - fees)/Price of Major Molecule.

Initial Major molecule token price = $1

Number of The fee is always swapped to USDC and sent to the major molecule fee wallet which stores all the revenue.

Price of Major molecule token at any point = (USD value of all assets within the major molecule +/- Margin Gains/Losses)/circulating supply of Major Molecule Tokens
```

### Fund management

Create a vault account address to manage funds per vault? Or manage all funds in molecule module account?

### Slippage

There should be well-formed slippage calculator considering oracle. (Consider GMX model)

We will introduce `external_liquidity` parameter to provide as less as slippage considering external exchanges.

One option is to introduce dynamically configurable slippage by governance of molecule token.

Volatility could be configured for individual assets, and slippage could be configured for pool itself.

For 4 asset pool (ETH, BTC, ATOM, USDC), different slippage could be executed based on following option.

Volatility of swapping assets + pool slippage + target weight change

TODO: should have exact Maths formula for calculating slippage.

### Asset weight recovery once it's broken

Asset weight recovery when it's broken is a challenging problem without incentive structure.

Following items are applied to keep the weight not broken

1. Slippage on swap when imbalanced
2. More fees on deposit when imbalanced

When weight is broken, users won't be using the Swap feature.
My opinion is for specific period when weight is lower than threshold, partial fees collected could be given to users based on the weight they recover.

Let’s assume there’s $100 fee collected and weight broken to 75%:25% level (edited)
If the pool enters weight rebalance period, users who swap JUNO to USDC gets fees.
I think this could be implemented something like funding fee on perpetual market.
People who break weight pay more fees
People who help weight balanced get small incentives when the weight is broken more than threshold
I think we don’t require a lot of funds to rebalance the weight
If we have $1K, we can do the swap operation repetitively in case fees are not enough and we will need to do ourselves.

## Endpoints

### Gov Proposals

- CreateVault(assets, ratios)
- UpdateVaultConfig(vaultId, targetWeights)

### Msg endpoints

- Deposit(assets) - calculate slippage based on weight change
- Swap(asset->target_asset) - calculate slippage based on weight change
- Withdraw(lp->target_asset) - calculate slippage based on weight change

### Query endpoints

- QueryVaults - Vault configs
- QueryVault(vaultId) - LP price, vault config, vault TVL, fees collected
- EstimatedSwapOutAmount(SwapCoin, outToken) -> amount
- EstimatedLPTokenAmountAfterDeposit(depositCoins) -> amount
- EstimatedWithdrawAmount(lpTokenAmount, outToken) -> outTokenAmount

### Epoch actions

#### Fee distribution

Fees are transferred to fee collector wallet and distributed per epoch to LPs and stakers per epoch.

LPs claim rewards through `commitment` module.
For stakers, will need to distribute through `incentive` module.

## Risk management

- Need to have a way to stop deposit/withdrawal of specific token within the vault?
- Need to have a way to stop overall deposit/withdrawal on the vault?
- When oracle has large difference between several platforms, it's not trustworthy, and Molecule token will need ask high fees so that it does not give more than the minimum price or just redirect the operation to normal Elys AMM.

## To consider

- Protect LP from volatile markets and market arbs
- Oracle difference % from sudden dump on third party market
- Weights are balanced by the incentive structure to hit close to the target weights
- Unit test to check final slippage on example values
- Low liquidity coin and high liquidity coin (low liquidity coin won't be able to use full oracle price)
- Once the user buys the major molecule TKN, it is automatically “ committed” in the commitment molecule to receive rewards
