package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	etypes "github.com/elys-network/elys/x/epochs/types"
	"github.com/elys-network/elys/x/incentive/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
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

// Find out active incentive params
func (k Keeper) GetProperIncentiveParam(ctx sdk.Context, epochIdentifier string) (bool, types.IncentiveInfo, types.IncentiveInfo) {
	// Fetch incentive params
	params := k.GetParams(ctx)

	// If we don't have enough params
	if len(params.StakeIncentives) < 1 || len(params.LpIncentives) < 1 {
		return false, types.IncentiveInfo{}, types.IncentiveInfo{}
	}

	// Current block timestamp
	timestamp := ctx.BlockTime().Unix()
	foundIncentive := false

	// Incentive params initialize
	stakeIncentive := params.StakeIncentives[0]
	lpIncentive := params.LpIncentives[0]

	// Consider epochIdentifier and start time
	// Consider epochNumber as well
	if stakeIncentive.EpochIdentifier != epochIdentifier || timestamp < stakeIncentive.StartTime.Unix() {
		return false, types.IncentiveInfo{}, types.IncentiveInfo{}
	}

	// Increase current epoch of Stake incentive param
	stakeIncentive.CurrentEpoch = stakeIncentive.CurrentEpoch + 1
	if stakeIncentive.CurrentEpoch == stakeIncentive.NumEpochs {
		params.StakeIncentives = params.StakeIncentives[1:]
	}

	// Increase current epoch of Lp incentive param
	lpIncentive.CurrentEpoch = lpIncentive.CurrentEpoch + 1
	if lpIncentive.CurrentEpoch == lpIncentive.NumEpochs {
		params.LpIncentives = params.LpIncentives[1:]
	}

	// Update params
	k.SetParams(ctx, params)

	// return found, stake, lp incentive params
	return foundIncentive, stakeIncentive, lpIncentive
}

// Calculate epoch counts per year to be used in APR calculation
func (k Keeper) CalculateEpochCountsPerYear(epochIdentifier string) int64 {
	switch epochIdentifier {
	case etypes.WeekEpochID:
		return ptypes.WeeksPerYear
	case etypes.DayEpochID:
		return ptypes.DaysPerYear
	case etypes.HourEpochID:
		return ptypes.HoursPerYear
	}

	return 0
}

// Update total commitment info
func (k Keeper) UpdateTotalCommitmentInfo(ctx sdk.Context) {
	// Fetch total staked Elys amount again
	k.tci.TotalElysBonded = k.stk.TotalBondedTokens(ctx)
	// Initialize with amount zero
	k.tci.TotalEdenEdenBoostCommitted = sdk.ZeroInt()
	// Initialize with amount zero
	k.tci.TotalFeesCollected = sdk.Coins{}
	// Initialize Lp tokens amount
	k.tci.TotalLpTokensCommitted = make(map[string]sdk.Int)

	// Collect gas fees collected
	fees := k.CollectGasFeesToIncentiveModule(ctx)
	// Calculate total fees - DEX revenus + Gas fees collected
	k.tci.TotalFeesCollected = k.tci.TotalFeesCollected.Add(fees...)

	// Iterate to calculate total Eden, Eden boost and Lp tokens committed
	k.cmk.IterateCommitments(ctx, func(commitments ctypes.Commitments) bool {
		committedEdenToken := commitments.GetCommittedAmountForDenom(ptypes.Eden)
		committedEdenBoostToken := commitments.GetCommittedAmountForDenom(ptypes.EdenB)

		k.tci.TotalEdenEdenBoostCommitted = k.tci.TotalEdenEdenBoostCommitted.Add(committedEdenToken).Add(committedEdenBoostToken)

		// Iterate to calcaulte total Lp tokens committed
		k.lpk.IterateLiquidityPools(ctx, func(l LiquidityPool) bool {
			committedLpToken := commitments.GetCommittedAmountForDenom(l.lpToken)
			k.tci.TotalLpTokensCommitted[l.lpToken] = k.tci.TotalLpTokensCommitted[l.lpToken].Add(committedLpToken)
			return false
		})
		return false
	})
}
