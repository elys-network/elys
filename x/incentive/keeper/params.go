package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx sdk.Context) (percent sdk.Dec) {
	k.paramstore.Get(ctx, types.ParamStoreKeyCommunityTax, &percent)
	return percent
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx sdk.Context) (enabled bool) {
	k.paramstore.Get(ctx, types.ParamStoreKeyWithdrawAddrEnabled, &enabled)
	return enabled
}

// GetDEXRewardPortionForLPs returns the dex revenue percent for Lps
func (k Keeper) GetDEXRewardPortionForLPs(ctx sdk.Context) (percent sdk.Dec) {
	k.paramstore.Get(ctx, types.ParamStoreKeyRewardPortionForLps, &percent)
	return percent
}

// GetDEXRewardPortionForStakers returns the dex revenue percent for Stakers
func (k Keeper) GetDEXRewardPortionForStakers(ctx sdk.Context) (percent sdk.Dec) {
	k.paramstore.Get(ctx, types.ParamStoreKeyRewardPortionForStakers, &percent)
	return percent
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

// UpdatePoolMultipliers updates pool multipliers through gov proposal
func (k Keeper) UpdatePoolMultipliers(ctx sdk.Context, poolMultipliers []types.PoolMultiplier) bool {
	if len(poolMultipliers) < 1 {
		return false
	}

	// Fetch incentive params
	params := k.GetParams(ctx)

	// Update pool multiplier
	for _, pm := range poolMultipliers {
		for i, p := range params.PoolInfos {
			// If we found matching poolId
			if p.PoolId == pm.PoolId {
				params.PoolInfos[i].Multiplier = pm.Multiplier
			}
		}
	}

	// Update parameter
	k.SetParams(ctx, params)

	return true
}

// Calculate epoch counts per year to be used in APR calculation
func (k Keeper) CalculateEpochCountsPerYear(ctx sdk.Context, epochIdentifier string) int64 {
	epochInfo, found := k.epochsKeeper.GetEpochInfo(ctx, epochIdentifier)
	epochSeconds := int64(epochInfo.Duration.Seconds())
	if !found || epochSeconds == 0 {
		return 0
	}

	// epoch min & max check
	if epochSeconds == 0 || epochSeconds > ptypes.SecondsPerYear {
		return 0
	}

	// returns num of epochs
	return ptypes.SecondsPerYear / epochSeconds
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
	k.tci.TotalLpTokensCommitted = make(map[string]sdk.Int)
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
