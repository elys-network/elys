package keeper

import (
	"fmt"
	"math"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	ctypes "github.com/elys-network/elys/x/commitment/types"
	etypes "github.com/elys-network/elys/x/epochs/types"
	"github.com/elys-network/elys/x/incentive/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace
		cmk        types.CommitmentKeeper
		stk        types.StakingKeeper
		tci        *types.TotalCommitmentInfo
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	ck types.CommitmentKeeper,
	sk types.StakingKeeper,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		cmk:        ck,
		stk:        sk,
		tci:        &types.TotalCommitmentInfo{},
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// Update total commitment info
func (k Keeper) UpdateTotalCommitmentInfo(ctx sdk.Context) {
	// Fetch total staked Elys amount again
	k.tci.TotalElysBonded = k.stk.TotalBondedTokens(ctx)
	// Initialize with amount zero
	k.tci.TotalCommitted = sdk.ZeroInt()

	// Iterate to calculate total Eden and Eden boost committed
	k.cmk.IterateCommitments(ctx, func(commitments ctypes.Commitments) bool {
		committedEdenToken := commitments.GetCommittedAmountForDenom(types.Eden)
		committedEdenBoostToken := commitments.GetCommittedAmountForDenom(types.EdenB)

		k.tci.TotalCommitted = k.tci.TotalCommitted.Add(committedEdenToken).Add(committedEdenBoostToken)

		return false
	})
}

// Calculate total share of staking
func (k Keeper) CalculateTotalShareOfStaking(amount sdk.Int) sdk.Dec {
	// Total statked = Elys staked + Eden Committed + Eden boost Committed
	totalStaked := k.tci.TotalElysBonded.Add(k.tci.TotalCommitted)

	// Share = Amount / Total Staked
	return sdk.NewDecFromInt(amount).QuoInt(totalStaked)
}

// Calculate the delegated amount
func (k Keeper) CalculateDelegatedAmount(ctx sdk.Context, delegator string) sdk.Int {
	// Derivate bech32 based delegator address
	delAdr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdk.ZeroInt()
	}

	// Get elys delegation for creator address
	delegatedAmt := sdk.ZeroDec()

	// Get all delegations
	delegations := k.stk.GetDelegatorDelegations(ctx, delAdr, math.MaxUint16)
	for _, del := range delegations {
		// Get validator address
		valAddr := del.GetValidatorAddr()
		// Get validator
		val := k.stk.Validator(ctx, valAddr)

		shares := del.GetShares()
		tokens := val.TokensFromSharesTruncated(shares)
		delegatedAmt = delegatedAmt.Add(tokens)
	}

	return delegatedAmt.TruncateInt()
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

// Update uncommitted token amount
// Called back through epoch hook
func (k Keeper) UpdateUncommittedTokens(ctx sdk.Context, epochIdentifier string, stakeIncentive types.IncentiveInfo, lpIncentive types.IncentiveInfo) {
	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx)

	// Calculate eden amount per epoch
	edenAmountPerEpoch := stakeIncentive.Amount.Quo(sdk.NewInt(stakeIncentive.NumEpochs))
	edenBoostAPR := stakeIncentive.EdenBoostApr

	// Iterate all delegations for the specified delegator
	// Process to increase uncomitted token amount of Eden & Eden boost
	k.cmk.IterateCommitments(
		ctx, func(commitments ctypes.Commitments) bool {
			// Commitment owner
			creator := commitments.Creator

			// Calculate delegated amount per delegator
			delegatedAmt := k.CalculateDelegatedAmount(ctx, creator)

			// Calculate new uncommitted Eden tokens for LP, staker, and Eden token holders
			newUncommittedEdenTokens := k.CalculateNewUncommittedEdenTokens(ctx, delegatedAmt, commitments, edenAmountPerEpoch)

			// Calculate new uncommitted Eden-Boost tokens for staker and Eden token holders
			newUncommittedEdenBoostTokens := k.CalculateNewUncommittedEdenBoostTokens(ctx, delegatedAmt, commitments, epochIdentifier, edenBoostAPR)

			// Update Commitments with new uncommitted token amounts
			k.UpdateCommitments(ctx, creator, &commitments, newUncommittedEdenTokens, newUncommittedEdenBoostTokens)

			return false
		},
	)
}

// Calculate new Eden token amounts based on the given conditions and user's current uncommitted token balance
func (k Keeper) CalculateNewUncommittedEdenTokens(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, edenAmountPerEpoch sdk.Int) sdk.Int {
	// Get LP commitments - Skip for now
	edenCommittedByLP := sdk.ZeroInt()

	// Get eden commitments and eden boost commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(types.Eden)
	edenBoostCommitted := commitments.GetCommittedAmountForDenom(types.EdenB)

	// compute eden reward based on above and param factors for each
	totalEdenCommittedByStake := delegatedAmt.Add(edenCommitted).Add(edenBoostCommitted).Add(edenCommittedByLP)
	stakeShare := k.CalculateTotalShareOfStaking(totalEdenCommittedByStake)

	// Calculate newly creating eden amount by its share
	newEdenAllocated := stakeShare.MulInt(edenAmountPerEpoch)

	return newEdenAllocated.TruncateInt()
}

// Calculate epoch counts per year to be used in APR calculation
func (k Keeper) CalculateEpochCountsPerYear(epochIdentifier string) int64 {
	switch epochIdentifier {
	case etypes.WeekEpochID:
		return types.WeeksPerYear
	case etypes.DayEpochID:
		return types.DaysPerYear
	case etypes.HourEpochID:
		return types.HoursPerYear
	}

	return 0
}

// Calculate new Eden-Boost token amounts based on the given conditions and user's current uncommitted token balance
func (k Keeper) CalculateNewUncommittedEdenBoostTokens(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, epochIdentifier string, edenBoostAPR int64) sdk.Int {
	// Get eden commitments
	edenCommitted := commitments.GetCommittedAmountForDenom(types.Eden)

	// Gompute eden reward based on above and param factors for each
	totalEden := delegatedAmt.Add(edenCommitted)

	// Calculate edenBoostAPR % APR for eden boost
	epochNumsPerYear := k.CalculateEpochCountsPerYear(epochIdentifier)

	return totalEden.Quo(sdk.NewInt(epochNumsPerYear)).Quo(sdk.NewInt(100)).Mul(sdk.NewInt(edenBoostAPR))
}

func (k Keeper) UpdateCommitments(ctx sdk.Context, creator string, commitments *ctypes.Commitments, newUncommittedEdenTokens sdk.Int, newUncommittedEdenBoostTokens sdk.Int) {
	// Update uncommitted Eden and Eden-Boost token balances in the Commitments structure
	k.UpdateEdenTokens(commitments, newUncommittedEdenTokens)
	k.UpdateEdenBoostTokens(commitments, newUncommittedEdenBoostTokens)

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, *commitments)
}

// Update the uncommitted Eden token balance
func (k Keeper) UpdateEdenTokens(commitments *ctypes.Commitments, new_uncommitted_eden_tokens sdk.Int) {
	uncommittedEden, found := commitments.GetUncommittedTokensForDenom(types.Eden)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedTokens = append(uncommittedTokens, &ctypes.UncommittedTokens{
			Denom:  types.Eden,
			Amount: new_uncommitted_eden_tokens,
		})
		commitments.UncommittedTokens = uncommittedTokens
	} else {
		uncommittedEden.Amount = uncommittedEden.Amount.Add(new_uncommitted_eden_tokens)
	}
}

// Update the uncommitted Eden-Boost token balance
func (k Keeper) UpdateEdenBoostTokens(commitments *ctypes.Commitments, new_uncommitted_eden_boost_tokens sdk.Int) {
	uncommittedEdenBoost, found := commitments.GetUncommittedTokensForDenom(types.Eden)
	if !found {
		uncommittedTokens := commitments.GetUncommittedTokens()
		uncommittedTokens = append(uncommittedTokens, &ctypes.UncommittedTokens{
			Denom:  types.EdenB,
			Amount: new_uncommitted_eden_boost_tokens,
		})
		commitments.UncommittedTokens = uncommittedTokens
	} else {
		uncommittedEdenBoost.Amount = uncommittedEdenBoost.Amount.Add(new_uncommitted_eden_boost_tokens)
	}
}
