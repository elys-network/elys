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

Generates an address to be used in collecting DEX revenue and perpetual revenue from the specified pool. We will have unique reward wallet address per pool and incentive module will access these wallets and distribute them to the specified LPs.

## Elys token to USDC conversion

Gas fees will be collected in Elys. It should be converted into USDC before being distributed. Other DEX fees will be collected in USDC. Meaning incentive module doesn't need to work on converting DEX fees collected. To convert Elys to USDC, we can use the following function to execute swap operation and have USDC converted.

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

1. Inflationary rewards distribution
   We need to iterate LP pools and then calculate the rewarding amount for LPs of the specified pool.
   We use proxy TVL rather than pure TVL for calculating each pool share. Given the Eden amount per epoch, we use multiplier of each pool. We can calculate proxy TVL by using multiplier and then can calculate each pool share.

Each pool has unique LP token and it is committed to commitment module. We can calculate the total committed LP token amount and amount LP token committed by specific liquidity provider. We can then calculate his share to the pool and then finally we calculate his Eden allocation during the epoch.

```go
// Proxy TVL
// Multiplier on each liquidity pool
// We have 3 pools of 20, 30, 40 TVL
// We have mulitplier of 0.3, 0.5, 1.0
// Proxy TVL = 20*0.3+30*0.5+40*1.0
Rewards = RewardsAmountOfPoolPerEpoch * UserLPCommitment / TotalLPCommitment
```

2. Non-inflationary rewards distribution
   We need to iterate LP reward pools with positive balance and then calculate the rewarding amount for LPs of the specified pool. We should calculate the LPs share of the specified pool.

```go
Rewards = l.rewardAmount * UserLPCommitment / TotalLPCommitmentOfPool
```

Rewards for LPs should be combined with Eden, Eden boost and DEX/perpetual rewards.

Note: To update in rewards distribution - users can run a bot to instantly put liquidity just before reward distribution and take it out. This happens because of zero lockup time and epoch based reward distribution. One of lockup or epoch based distribution modified. This is common problem on reward distribution per epoch.
The easiest way to solve this issue is by timestamping whenever an LP add or withdraw liquidity and only distribute rewards to the LP if the timestamp is older than the epoch

## Reward distribution logic for stakers

1. Inflationary rewards distribution
   Given the amount of Eden per epoch, we calculate the amount of Eden and Eden boost for stakers.

```go
    // Calculate total share of staking considering Eden committed, Eden boost committed and Elys staked.
	stakeShare := k.CalcTotalShareOfStaking(totalEdenCommittedByStake)
	newEdenAllocated := stakeShare.MulInt(edenAmountPerEpoch)
```

2. Non-inflationary rewards distribution
   35% of all DEX revenue is distributed to the stakers according to their staking share.

```go
    dexRewards := stakeShare.Mul(dexRevenueAmtForStakers).TruncateInt()
```

## Withdraw rewards

1. Withdraw rewards of LPs and stakers
2. Withdraw commissions of validators.
