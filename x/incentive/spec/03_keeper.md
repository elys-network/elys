<!--
order: 3
-->

# Keepers

## Reward pool address for LPs

There's constant reward pool address per liquidity pool id which will be distributed per epoch.

```go
    func GetLPRewardsPoolAddress(poolId uint64) {
        return authtypes.NewModuleAddress(fmt.Sprintf("lp_rewards_pool_%d", poolId))
    }
```

DEX revenue and margin revenue to be distributed to specific pool liquidity providers is sent to the address and incentive module distribute the those tokens to relevant pool LPs.

## Elys token to USDC conversion

Before reward distribution, it will be needed to check if it has positive Elys balance and in this case, execute swap operation.

```go
ammKeeper.SwapExactAmountIn(goCtx, &types.MsgSwapExactAmountIn{
    Creator: rewardPoolAddress,
	TokenIn: elysAmount,
	TokenOutMinAmount: ElysPrice * elysAmount * 0.9,
	SwapRoutePoolIds: []uint64{elysPoolId}, // elysPoolId to be configured on incentive module params
	SwapRouteDenoms   []string{},
})
```

## Reward distribution logic for LPs

Reward distribution for LPs will need to iterate LP reward pools with positive balance.
And configure total rewards for LPs for pool.

```go
Rewards = EpochTotalRewards * UserLPCommitment / TotalLPCommitment
```

Rewards for LPs should be combined with Eden, Eden boost and DEX/margin rewards.

TODO: Problem in this logic - Users can run a bot to instantly put liquidity just before reward distribution and take it out. This happens because of zero lockup time and epoch based reward distribution. One of lockup or epoch based distribution modified. This is common problem on reward distribution per epoch.
Unless there's a way to track if the commitment was done past epoch and not unbonded so far.

## Reward distribution logic for stakers

This mechanism should be similar to Cosmos SDK's distribution module does except for epoch based operation.
