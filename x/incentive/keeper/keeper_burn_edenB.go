package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Burn EdenBoost from Elys unstaked
func (k Keeper) BurnEdenBFromElysUnstaking(ctx sdk.Context, delegator sdk.AccAddress) {
	delAddr := delegator.String()
	// Get commitments
	commitments := k.cmk.GetCommitments(ctx, delAddr)

	// Get previous amount
	prevElysStaked := k.GetElysStaked(ctx, delAddr)
	if prevElysStaked.IsZero() {
		return
	}

	// Calculate current delegated amount of delegator
	delAmount := k.CalcDelegationAmount(ctx, delAddr)

	// If not unstaked,
	// should return nil otherwise it will break staking module
	if delAmount.GTE(prevElysStaked) {
		return
	}

	// TODO: might need to claim all rewards before burn operation to properly burn EdenB including unclaimed

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBClaimed)

	// Unstaked
	unstakedElys := prevElysStaked.Sub(delAmount)

	unstakedElysDec := sdk.NewDecFromInt(unstakedElys)
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(prevElysStaked))
	totalEdenBDec := sdk.NewDecFromInt(totalEdenB)
	edenBToBurn := sdk.ZeroDec()
	if edenCommittedAndElysStakedDec.GT(sdk.ZeroDec()) {
		edenBToBurn = unstakedElysDec.Quo(edenCommittedAndElysStakedDec).Mul(totalEdenBDec)
	}
	// Burn EdenB in commitment module
	commitment, err := k.cmk.BurnEdenBoost(ctx, delAddr, ptypes.EdenB, edenBToBurn.TruncateInt())
	if err != nil {
		k.Logger(ctx).Error("EdenB burn failure", err)
	} else {
		k.cmk.SetCommitments(ctx, commitment)
	}
}

// Burn EdenBoost from Eden unclaimed
func (k Keeper) BurnEdenBFromEdenUncommitted(ctx sdk.Context, delegator string, uncommitAmt math.Int) error {
	// Get elys staked amount
	elysStaked := k.GetElysStaked(ctx, delegator)
	if elysStaked.IsZero() {
		return nil
	}

	commitments := k.cmk.GetCommitments(ctx, delegator)
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBClaimed)

	unclaimedAmtDec := sdk.NewDecFromInt(uncommitAmt)
	// This formula shud be applied before eden uncommitted or elys staked is removed from eden committed amount and elys staked amount respectively
	// So add uncommitted amount to committed eden bucket in calculation.
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(elysStaked).Add(uncommitAmt))
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
