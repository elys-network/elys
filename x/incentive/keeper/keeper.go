package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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
func (k Keeper) UpdateTotalCommitmentInfo(ctx sdk.Context, epochIdentifier string) {
	// Fetch total staked Elys amount again
	k.tci.TotalElysBonded = k.stk.TotalBondedTokens(ctx)
	// Initialize with amount zero
	k.tci.TotalEdenCommitted = sdk.ZeroInt()
	k.tci.TotalEdenBoostCommitted = sdk.ZeroInt()

	// Iterate to calculate total Eden and Eden boost committed
	k.cmk.IterateCommitments(ctx, epochIdentifier, func(commitments ctypes.Commitments) bool {
		committedTokens := commitments.CommittedTokens
		// Sum Eden and Eden boost committed
		for _, c := range committedTokens {
			// Eden
			if c.Denom == types.Eden {
				k.tci.TotalEdenCommitted = k.tci.TotalEdenCommitted.Add(c.Amount)
			}

			// Eden boost
			if c.Denom == types.EdenB {
				k.tci.TotalEdenBoostCommitted = k.tci.TotalEdenCommitted.Add(c.Amount)
			}
		}

		return false
	})
}

// Calculate total share of staking
func (k Keeper) CalculateTotalShareOfStaking(amount sdk.Int) sdk.Dec {
	// Total statked = Elys staked + Eden Committed + Eden boost Committed
	totalStaked := k.tci.TotalElysBonded.Add(k.tci.TotalEdenCommitted).Add(k.tci.TotalEdenBoostCommitted)

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
	delegatedAmt := sdk.ZeroInt()

	// Iterate all delegations for the specified delegator
	k.stk.IterateDelegations(
		ctx, delAdr,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			// Get validator address
			valAddr := del.GetValidatorAddr()
			// Get validator
			val := k.stk.Validator(ctx, valAddr)

			/****************************************************************************************/
			// --------------------(Cosmos SDK staking module implementation)------------------------
			// del.Shares = val.GetDelegatorShares().MulInt(amt).QuoInt(val.GetTokens())
			// amt = del.Shares * val.GetTokens / val.GetDelegatorShares()
			/****************************************************************************************/
			amt := del.GetShares().MulInt(val.GetTokens()).Quo(val.GetDelegatorShares()).TruncateInt()

			// Sum the individual delegation
			delegatedAmt.Add(amt)
			return false
		},
	)

	return delegatedAmt
}

// Find out active incentive params
func (k Keeper) FindProperIncentiveParm(ctx sdk.Context, epochIdentifier string) (bool, types.IncentiveInfo, types.IncentiveInfo) {
	// Incentive params initialize
	stakeIncentive := types.IncentiveInfo{}
	lpIncentive := types.IncentiveInfo{}

	// Current block timestamp
	timestamp := ctx.BlockTime().Unix()
	foundIncentive := false

	// Fetch incentive params
	params := k.GetParams(ctx)

	// Find approporiate Incentive Info
	// Consider epochIdentifier and start time
	// Consider epochNumber as well
	for _, ii := range params.StakeIncentives {
		if ii.EpochIdentifier == epochIdentifier && timestamp >= ii.StartTime.Unix() {
			stakeIncentive = ii
			foundIncentive = true
			break
		}
	}

	// Find approporiate Incentive Info
	for _, ii := range params.LPIncentives {
		if ii.EpochIdentifier == epochIdentifier && timestamp >= ii.StartTime.Unix() {
			lpIncentive = ii
		}
	}

	// return found, stake, lp incentive params
	return foundIncentive, stakeIncentive, lpIncentive
}

// Update uncommitted token amount
// Called back through epoch hook
func (k Keeper) UpdateUncommittedTokens(ctx sdk.Context, epochIdentifier string) {
	// Recalculate total committed info
	k.UpdateTotalCommitmentInfo(ctx, epochIdentifier)

	// Iterate all delegations for the specified delegator
	// Process to increase uncomitted token amount of Eden & Eden boost
	k.cmk.IterateCommitments(
		ctx, epochIdentifier,
		func(commitments ctypes.Commitments) bool {
			// Find out active incentive params
			foundIncentive, stakeIncentive, lpIncentive := k.FindProperIncentiveParm(ctx, epochIdentifier)

			// If we don't have incentive params ready
			if !foundIncentive {
				return true
			}

			// Commitment owner
			creator := commitments.Creator

			// Calculate delegated amount per delegator
			delegatedAmt := k.CalculateDelegatedAmount(ctx, creator)

			// Calculate new uncommitted Eden tokens for LP, staker, and Eden token holders
			new_uncommitted_eden_tokens := k.CalculateNewUncommittedEdenTokens(ctx, delegatedAmt, commitments, stakeIncentive, lpIncentive)

			// Calculate new uncommitted Eden-Boost tokens for staker and Eden token holders
			new_uncommitted_eden_boost_tokens := k.CalculateNewUncommittedEdenBoostTokens(ctx, delegatedAmt, commitments, epochIdentifier, stakeIncentive, lpIncentive)

			// Update Commitments with new uncommitted token amounts
			k.UpdateCommitments(ctx, creator, new_uncommitted_eden_tokens, new_uncommitted_eden_boost_tokens)

			return false
		},
	)
}

// Calculate new Eden token amounts based on the given conditions and user's current uncommitted token balance
func (k Keeper) CalculateNewUncommittedEdenTokens(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, stakeIncentive types.IncentiveInfo, lpIncentive types.IncentiveInfo) sdk.Int {
	// Get LP commitments - Skip for now
	edenCommittedByLP := sdk.ZeroInt()

	// Get eden commitments and eden boost commitments
	edenCommitted := sdk.ZeroInt()
	edenBoostCommitted := sdk.ZeroInt()
	for _, c := range commitments.CommittedTokens {
		// Eden committed
		if c.Denom == types.Eden {
			edenCommitted = c.Amount
		}

		// Eden boost committed
		if c.Denom == types.EdenB {
			edenBoostCommitted = c.Amount
		}
	}

	// compute eden reward based on above and param factors for each
	totalEdenCommittedByStake := delegatedAmt.Add(edenCommitted).Add(edenBoostCommitted).Add(edenCommittedByLP)
	stakeShare := k.CalculateTotalShareOfStaking(totalEdenCommittedByStake)

	// Calculate eden amount per epoch
	edenAmountPerEpoch := stakeIncentive.Amount.Quo(sdk.NewInt(stakeIncentive.NumEpochs))

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
func (k Keeper) CalculateNewUncommittedEdenBoostTokens(ctx sdk.Context, delegatedAmt sdk.Int, commitments ctypes.Commitments, epochIdentifier string, stakeIncentive types.IncentiveInfo, lpIncentive types.IncentiveInfo) sdk.Int {
	// Get eden commitments
	edenCommitted := sdk.ZeroInt()
	for _, c := range commitments.CommittedTokens {
		// Eden committed
		if c.Denom == types.Eden {
			edenCommitted = c.Amount
		}
	}

	// Gompute eden reward based on above and param factors for each
	totalEden := delegatedAmt.Add(edenCommitted)

	// Calculate 100% APR for eden boost
	epochNumsPerYear := k.CalculateEpochCountsPerYear(epochIdentifier)

	return totalEden.Quo(sdk.NewInt(epochNumsPerYear))
}

func (k Keeper) UpdateCommitments(ctx sdk.Context, creator string, new_uncommitted_eden_tokens sdk.Int, new_uncommitted_eden_boost_tokens sdk.Int) {
	commitments, _ := k.cmk.GetCommitments(ctx, creator)

	// Update uncommitted Eden and Eden-Boost token balances in the Commitments structure
	k.UpdateEdenTokens(&commitments, new_uncommitted_eden_tokens)
	k.UpdateEdenBoostTokens(&commitments, new_uncommitted_eden_boost_tokens)

	// Save the updated Commitments
	k.cmk.SetCommitments(ctx, commitments)
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
