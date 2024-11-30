package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Burn EdenBoost from Elys unstaked
func (k Keeper) BurnEdenBFromElysUnstaking(ctx sdk.Context, delegator sdk.AccAddress) error {
	// Get commitments
	commitments := k.commKeeper.GetCommitments(ctx, delegator)

	// Get previous amount
	prevElysStaked := k.GetElysStaked(ctx, delegator)
	if prevElysStaked.Amount.IsZero() {
		return nil
	}

	// Calculate current delegated amount of delegator
	delAmount := k.CalcDelegationAmount(ctx, delegator)

	// If not unstaked,
	// should return nil otherwise it will break staking module
	if delAmount.GTE(prevElysStaked.Amount) {
		return nil
	}

	_, err := k.WithdrawAllRewards(ctx, &types.MsgWithdrawAllRewards{DelegatorAddress: delegator.String()})
	if err != nil {
		return err
	}

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)
	totalEdenB := edenBCommitted.Add(edenBClaimed)

	// Unstaked
	unstakedElys := prevElysStaked.Amount.Sub(delAmount)

	unstakedElysDec := math.LegacyNewDecFromInt(unstakedElys)
	edenCommittedAndElysStakedDec := math.LegacyNewDecFromInt(edenCommitted.Add(prevElysStaked.Amount))
	edenBToBurn := math.LegacyZeroDec()
	if edenCommittedAndElysStakedDec.GT(math.LegacyZeroDec()) {
		edenBToBurn = unstakedElysDec.Quo(edenCommittedAndElysStakedDec).MulInt(totalEdenB)
	}
	if edenBToBurn.IsZero() {
		return nil
	}

	// Burn EdenB in commitment module
	err = k.commKeeper.BurnEdenBoost(ctx, delegator, ptypes.EdenB, edenBToBurn.TruncateInt())
	if err != nil {
		return err
	}
	return nil
}

// Burn EdenBoost from Eden unclaimed
func (k Keeper) BurnEdenBFromEdenUncommitted(ctx sdk.Context, delegator sdk.AccAddress, uncommitAmt math.Int) error {
	_, err := k.WithdrawAllRewards(ctx, &types.MsgWithdrawAllRewards{DelegatorAddress: delegator.String()})
	if err != nil {
		return err
	}

	// Get elys staked amount
	elysStaked := k.GetElysStaked(ctx, delegator)
	commitments := k.commKeeper.GetCommitments(ctx, delegator)
	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	// Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBClaimed := commitments.GetClaimedForDenom(ptypes.EdenB)
	totalEdenB := edenBCommitted.Add(edenBClaimed)

	unclaimedAmtDec := math.LegacyNewDecFromInt(uncommitAmt)
	// This formula should be applied before eden uncommitted or elys staked is removed from eden committed amount and elys staked amount respectively
	// So add uncommitted amount to committed eden bucket in calculation.
	edenCommittedAndElysStakedDec := math.LegacyNewDecFromInt(edenCommitted.Add(elysStaked.Amount).Add(uncommitAmt))
	if edenCommittedAndElysStakedDec.IsZero() {
		return nil
	}

	edenBToBurn := math.LegacyZeroDec()
	if edenCommittedAndElysStakedDec.GT(math.LegacyZeroDec()) {
		edenBToBurn = unclaimedAmtDec.Quo(edenCommittedAndElysStakedDec).MulInt(totalEdenB)
	}
	if edenBToBurn.IsZero() {
		return nil
	}

	// Burn EdenB in commitment module
	err = k.commKeeper.BurnEdenBoost(ctx, delegator, ptypes.EdenB, edenBToBurn.TruncateInt())
	return err
}
