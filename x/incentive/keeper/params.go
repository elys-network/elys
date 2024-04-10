package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get([]byte(types.ParamsKey))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.ParamsKey), b)
}

// GetDEXRewardPortionForLPs returns the dex revenue percent for Lps
func (k Keeper) GetDEXRewardPortionForLPs(ctx sdk.Context) (percent sdk.Dec) {
	return k.GetParams(ctx).RewardPortionForLps
}

// GetDEXRewardPortionForStakers returns the dex revenue percent for Stakers
func (k Keeper) GetDEXRewardPortionForStakers(ctx sdk.Context) (percent sdk.Dec) {
	return k.GetParams(ctx).RewardPortionForStakers
}

// GetPoolInfo
func (k Keeper) GetPoolInfo(ctx sdk.Context, poolId uint64) (types.PoolInfo, bool) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	poolInfos := params.PoolInfos
	for _, ps := range poolInfos {
		if ps.PoolId == poolId {
			return ps, true
		}
	}

	return types.PoolInfo{}, false
}

// SetPoolInfo
func (k Keeper) SetPoolInfo(ctx sdk.Context, poolId uint64, poolInfo types.PoolInfo) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)

	poolInfos := params.PoolInfos
	for i, ps := range poolInfos {
		if ps.PoolId == poolId {
			params.PoolInfos[i] = poolInfo
			k.SetParams(ctx, params)

			return true
		}
	}

	return false
}

// InitPoolParams: creates a poolInfo at the time of pool creation.
func (k Keeper) InitPoolParams(ctx sdk.Context, poolId uint64) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)
	poolInfos := params.PoolInfos

	for _, ps := range poolInfos {
		if ps.PoolId == poolId {
			return true
		}
	}

	// Initiate a new pool info
	poolInfo := types.PoolInfo{
		// reward amount
		PoolId: poolId,
		// reward wallet address
		RewardWallet: ammtypes.NewPoolRevenueAddress(poolId).String(),
		// multiplier for lp rewards
		Multiplier: sdk.NewDec(1),
		// Number of blocks since creation
		NumBlocks: sdk.NewInt(1),
		// Total dex rewards given since creation
		DexRewardAmountGiven: sdk.ZeroDec(),
		// Total eden rewards given since creation
		EdenRewardAmountGiven: sdk.ZeroInt(),
	}

	// Update pool information
	params.PoolInfos = append(params.PoolInfos, poolInfo)
	k.SetParams(ctx, params)

	return true
}

// InitStableStakePoolMultiplier: create a stable stake pool information responding to the pool creation.
func (k Keeper) InitStableStakePoolParams(ctx sdk.Context, poolId uint64) bool {
	// Fetch incentive params
	params := k.GetParams(ctx)
	poolInfos := params.PoolInfos

	for _, ps := range poolInfos {
		if ps.PoolId == poolId {
			return true
		}
	}

	// Initiate a new pool info
	poolInfo := types.PoolInfo{
		// reward amount
		PoolId: poolId,
		// reward wallet address
		RewardWallet: stabletypes.PoolAddress().String(),
		// multiplier for lp rewards
		Multiplier: sdk.NewDec(1),
		// Number of blocks since creation
		NumBlocks: sdk.NewInt(1),
		// Total dex rewards given since creation
		DexRewardAmountGiven: sdk.ZeroDec(),
		// Total eden rewards given since creation
		EdenRewardAmountGiven: sdk.ZeroInt(),
	}

	// Update pool information
	params.PoolInfos = append(params.PoolInfos, poolInfo)
	k.SetParams(ctx, params)

	return true
}

// Update total commitment info
func (k Keeper) UpdateTotalCommitmentInfo(ctx sdk.Context, baseCurrency string) {
	// Fetch total staked Elys amount again
	k.tci.TotalElysBonded = k.stk.TotalBondedTokens(ctx)
	// Initialize with amount zero
	k.tci.TotalEdenEdenBoostCommitted = sdk.ZeroInt()
	// Initialize with amount zero
	k.tci.TotalFeesCollected = sdk.Coins{}
	// Initialize Lp tokens amount
	k.tci.TotalLpTokensCommitted = make(map[string]math.Int)
	// Reinitialize Pool revenue tracker
	k.tci.PoolRevenueTrack = make(map[string]sdk.Dec)

	// Collect gas fees collected
	fees := k.CollectGasFeesToIncentiveModule(ctx, baseCurrency)

	// Calculate total fees - Gas fees collected
	k.tci.TotalFeesCollected = k.tci.TotalFeesCollected.Add(fees...)

	// Iterate to calculate total Eden, Eden boost and Lp tokens committed
	k.cmk.IterateCommitments(ctx, func(commitments ctypes.Commitments) bool {
		committedEdenToken := commitments.GetCommittedAmountForDenom(ptypes.Eden)
		committedEdenBoostToken := commitments.GetCommittedAmountForDenom(ptypes.EdenB)

		k.tci.TotalEdenEdenBoostCommitted = k.tci.TotalEdenEdenBoostCommitted.Add(committedEdenToken).Add(committedEdenBoostToken)

		// Iterate to calculate total Lp tokens committed
		k.amm.IterateLiquidityPools(ctx, func(p ammtypes.Pool) bool {
			lpToken := ammtypes.GetPoolShareDenom(p.GetPoolId())

			committedLpToken := commitments.GetCommittedAmountForDenom(lpToken)
			amt, ok := k.tci.TotalLpTokensCommitted[lpToken]
			if !ok {
				k.tci.TotalLpTokensCommitted[lpToken] = committedLpToken
			} else {
				k.tci.TotalLpTokensCommitted[lpToken] = amt.Add(committedLpToken)
			}
			return false
		})

		// handle stable stake pool lp token
		lpStableStakeDenom := stabletypes.GetShareDenom()
		committedLpToken := commitments.GetCommittedAmountForDenom(lpStableStakeDenom)
		amt, ok := k.tci.TotalLpTokensCommitted[lpStableStakeDenom]
		if !ok {
			k.tci.TotalLpTokensCommitted[lpStableStakeDenom] = committedLpToken
		} else {
			k.tci.TotalLpTokensCommitted[lpStableStakeDenom] = amt.Add(committedLpToken)
		}
		return false
	})
}
