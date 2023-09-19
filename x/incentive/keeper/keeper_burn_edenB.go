package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Burn EdenBoost from Elys unstaked
func (k Keeper) BurnEdenBFromElysUnstaking(ctx sdk.Context, delegator sdk.AccAddress) error {
	delAddr := delegator.String()
	// Get commitments
	commitments, found := k.cmk.GetCommitments(ctx, delAddr)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address %s", delegator.String())
	}

	// Get previous amount
	prevElysStaked, found := k.GetElysStaked(ctx, delAddr)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address %s", delegator.String())
	}

	// Calculate current delegated amount of delegator
	delegatedAmt := k.CalculateDelegatedAmount(ctx, delAddr)

	// If not unstaked,
	if delegatedAmt.GTE(prevElysStaked.Amount) {
		return nil
	}

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	//Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBUncommitted := commitments.GetUncommittedAmountForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBUncommitted)

	// Unstaked
	unstakedElys := prevElysStaked.Amount.Sub(delegatedAmt)

	edenBToBurn := unstakedElys.Quo(edenCommitted.Add(prevElysStaked.Amount)).Mul(totalEdenB)

	// Burn EdenB ( Deduction EdenB in commitment module)
	_, err := k.cmk.DeductCommitments(ctx, delAddr, ptypes.EdenB, edenBToBurn)
	return err
}

// Burn EdenBoost from Eden uncommitted
func (k Keeper) BurnEdenBFromEdenUncommitted(ctx sdk.Context, delegator string, uncommittedAmt sdk.Int) error {
	// Get elys staked amount
	elysStaked, found := k.GetElysStaked(ctx, delegator)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address %s", delegator)
	}

	commitments, found := k.cmk.GetCommitments(ctx, delegator)
	if !found {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid address %s", delegator)
	}

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	//Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBUncommitted := commitments.GetUncommittedAmountForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBUncommitted)

	edenBToBurn := uncommittedAmt.Quo(edenCommitted.Add(elysStaked.Amount)).Mul(totalEdenB)

	// Burn EdenB ( Deduction EdenB in commitment module)
	_, err := k.cmk.DeductCommitments(ctx, delegator, ptypes.EdenB, edenBToBurn)
	return err
}
