package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Burn EdenBoost from Elys unstaked
func (k Keeper) BurnEdenBFromElysUnstaking(ctx sdk.Context, delegator sdk.AccAddress) error {
	delAddr := delegator.String()
	// Get commitments
	commitments, found := k.cmk.GetCommitments(ctx, delAddr)
	// should return nil otherwise it will break staking module
	if !found {
		return nil
	}

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
	edenBUncommitted := commitments.GetUncommittedAmountForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBUncommitted)

	// Unstaked
	unstakedElys := prevElysStaked.Amount.Sub(delegatedAmt)

	unstakedElysDec := sdk.NewDecFromInt(unstakedElys)
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(prevElysStaked.Amount))
	totalEdenBDec := sdk.NewDecFromInt(totalEdenB)
	edenBToBurn := unstakedElysDec.Quo(edenCommittedAndElysStakedDec).Mul(totalEdenBDec)

	// Burn EdenB ( Deduction EdenB in commitment module)
	commitment, err := k.cmk.DeductCommitments(ctx, delAddr, ptypes.EdenB, edenBToBurn.TruncateInt())
	k.cmk.SetCommitments(ctx, commitment)

	return err
}

// Burn EdenBoost from Eden uncommitted
func (k Keeper) BurnEdenBFromEdenUncommitted(ctx sdk.Context, delegator string, uncommittedAmt sdk.Int) error {
	// Get elys staked amount
	elysStaked, found := k.GetElysStaked(ctx, delegator)
	if !found {
		return nil
	}

	commitments, found := k.cmk.GetCommitments(ctx, delegator)
	// should return nil otherwise it will break commitment module
	if !found {
		return nil
	}

	edenCommitted := commitments.GetCommittedAmountForDenom(ptypes.Eden)

	//Total EdenB amount
	edenBCommitted := commitments.GetCommittedAmountForDenom(ptypes.EdenB)
	edenBUncommitted := commitments.GetUncommittedAmountForDenom(ptypes.EdenB)

	// Total EdenB amount
	totalEdenB := edenBCommitted.Add(edenBUncommitted)

	uncommittedAmtDec := sdk.NewDecFromInt(uncommittedAmt)
	edenCommittedAndElysStakedDec := sdk.NewDecFromInt(edenCommitted.Add(elysStaked.Amount))
	totalEdenBDec := sdk.NewDecFromInt(totalEdenB)

	edenBToBurn := uncommittedAmtDec.Quo(edenCommittedAndElysStakedDec).Mul(totalEdenBDec)

	// Burn EdenB ( Deduction EdenB in commitment module)
	commitment, err := k.cmk.DeductCommitments(ctx, delegator, ptypes.EdenB, edenBToBurn.TruncateInt())
	k.cmk.SetCommitments(ctx, commitment)

	return err
}
