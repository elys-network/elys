package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Burn EdenBoost from Elys unstaked
func (k Keeper) BurnEdenBFromElysUnstaking(ctx sdk.Context, delegator sdk.AccAddress) {
	// Get commitments
	commitments := k.commKeeper.GetCommitments(ctx, delegator)

	// Get previous amount
	prevElysStaked := k.GetElysStaked(ctx, delegator)
	if prevElysStaked.Amount.IsZero() {
		return
	}

	// Calculate current delegated amount of delegator
	delAmount := k.CalcDelegationAmount(ctx, delegator)

	// If not unstaked,
	// should return nil otherwise it will break staking module
	if delAmount.GTE(prevElysStaked.Amount) {
		return
	}

	// TODO: Claim all delegation rewards

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)
	totalEdenB := edenBCommitted.Add(edenBClaimed)

	// Unstaked
	unstakedElys := prevElysStaked.Amount.Sub(delAmount)

	unstakedElysDec := sdk.NewDecFromInt(unstakedElys)
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(prevElysStaked.Amount))
	edenBToBurn := sdk.ZeroDec()
	if edenCommittedAndElysStakedDec.GT(sdk.ZeroDec()) {
		edenBToBurn = unstakedElysDec.Quo(edenCommittedAndElysStakedDec).MulInt(totalEdenB)
	}
	if edenBToBurn.IsZero() {
		return
	}

	// Burn EdenB in commitment module
	cacheCtx, write := ctx.CacheContext()
	err := k.commKeeper.BurnEdenBoost(cacheCtx, delegator, ptypes.EdenB, edenBToBurn.TruncateInt())
	if err != nil {
		k.Logger(ctx).Error("EdenB burn failure", err)
	} else {
		write()
	}
}

// Burn EdenBoost from Eden unclaimed
func (k Keeper) BurnEdenBFromEdenUncommitted(ctx sdk.Context, delegator sdk.AccAddress, uncommitAmt math.Int) error {
	// TODO: Claim all delegation rewards

	// Get elys staked amount
	elysStaked := k.GetElysStaked(ctx, delegator)
	commitments := k.commKeeper.GetCommitments(ctx, delegator)
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)
	totalEdenB := edenBCommitted.Add(edenBClaimed)

	unclaimedAmtDec := sdk.NewDecFromInt(uncommitAmt)
	// This formula shud be applied before eden uncommitted or elys staked is removed from eden committed amount and elys staked amount respectively
	// So add uncommitted amount to committed eden bucket in calculation.
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(elysStaked.Amount).Add(uncommitAmt))
	if edenCommittedAndElysStakedDec.IsZero() {
		return nil
	}

	edenBToBurn := sdk.ZeroDec()
	if edenCommittedAndElysStakedDec.GT(sdk.ZeroDec()) {
		edenBToBurn = unclaimedAmtDec.Quo(edenCommittedAndElysStakedDec).MulInt(totalEdenB)
	}
	if edenBToBurn.IsZero() {
		return nil
	}

	// Burn EdenB in commitment module
	err := k.commKeeper.BurnEdenBoost(ctx, delegator, ptypes.EdenB, edenBToBurn.TruncateInt())
	return err
}
