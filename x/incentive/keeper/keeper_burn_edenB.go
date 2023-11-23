package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Burn EdenBoost from Elys unstaked
func (k Keeper) BurnEdenBFromElysUnstaking(ctx sdk.Context, delegator sdk.AccAddress) error {
	delAddr := delegator.String()
	// Get commitments
	commitments := k.cmk.GetCommitments(ctx, delAddr)

	// Get previous amount
	prevElysStaked, found := k.GetElysStaked(ctx, delAddr)
	// should return nil otherwise it will break staking module
	if !found {
		return nil
	}

	// Calculate current delegated amount of delegator
	delegatedAmt := k.CalculateDelegatedAmount(ctx, delAddr)

	// If not unstaked,
	// should return nil otherwise it will break staking module
	if delegatedAmt.GTE(prevElysStaked.Amount) {
		return nil
	}

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	//Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBUnclaimed := commitments.GetRewardUnclaimedForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBUnclaimed).Add(edenBClaimed)

	// Unstaked
	unstakedElys := prevElysStaked.Amount.Sub(delegatedAmt)

	unstakedElysDec := sdk.NewDecFromInt(unstakedElys)
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(prevElysStaked.Amount))
	totalEdenBDec := sdk.NewDecFromInt(totalEdenB)
	edenBToBurn := sdk.ZeroDec()
	if edenCommittedAndElysStakedDec.GT(sdk.ZeroDec()) {
		edenBToBurn = unstakedElysDec.Quo(edenCommittedAndElysStakedDec).Mul(totalEdenBDec)
	}
	// Burn EdenB ( Deduction EdenB in commitment module)
	commitment, err := k.cmk.BurnEdenBoost(ctx, delAddr, ptypes.EdenB, edenBToBurn.TruncateInt())
	k.cmk.SetCommitments(ctx, commitment)

	return err
}

// Burn EdenBoost from Eden unclaimed
func (k Keeper) BurnEdenBFromEdenUncommitted(ctx sdk.Context, delegator string, uncommitAmt sdk.Int) error {
	// Get elys staked amount
	elysStaked, found := k.GetElysStaked(ctx, delegator)
	if !found {
		return nil
	}

	commitments := k.cmk.GetCommitments(ctx, delegator)
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBUnclaimed := commitments.GetRewardUnclaimedForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBUnclaimed).Add(edenBClaimed)

	unclaimedAmtDec := sdk.NewDecFromInt(uncommitAmt)
	// This formula shud be applied before eden uncommitted or elys staked is removed from eden committed amount and elys staked amount respectively
	// So add uncommitted amount to committed eden bucket in calculation.
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(elysStaked.Amount).Add(uncommitAmt))
	totalEdenBDec := sdk.NewDecFromInt(totalEdenB)

	edenBToBurn := sdk.ZeroDec()
	if edenCommittedAndElysStakedDec.GT(sdk.ZeroDec()) {
		edenBToBurn = unclaimedAmtDec.Quo(edenCommittedAndElysStakedDec).Mul(totalEdenBDec)
	}

	// Burn EdenB ( Deduction EdenB in commitment module)
	commitment, err := k.cmk.BurnEdenBoost(ctx, delegator, ptypes.EdenB, edenBToBurn.TruncateInt())
	k.cmk.SetCommitments(ctx, commitment)

	return err
}
