package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Increase uncommitted token amount for the corresponding validator
func (k Keeper) UpdateTokensForValidator(ctx sdk.Context, validator string, new_uncommitted_eden_tokens sdk.Int, dexRewards sdk.Dec) {
	commitments, bfound := k.cmk.GetCommitments(ctx, validator)
	if !bfound {
		return
	}

	// Update Eden amount
	k.UpdateTokensCommitment(&commitments, new_uncommitted_eden_tokens, ptypes.Eden)

	// Update Elys amount
	k.UpdateTokensCommitment(&commitments, dexRewards.TruncateInt(), ptypes.Elys)

	// Update commmitment
	k.cmk.SetCommitments(ctx, commitments)
}

// Give commissions to validators
func (k Keeper) GiveCommissionToValidators(ctx sdk.Context, delegator string, totalDelegationAmt sdk.Int, newUncommittedAmt sdk.Int, dexRewards sdk.Dec) (sdk.Int, sdk.Int) {
	delAdr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdk.ZeroInt(), sdk.ZeroInt()
	}

	// If there is no delegation, (not elys staker)
	if totalDelegationAmt.LTE(sdk.ZeroInt()) {
		return sdk.ZeroInt(), sdk.ZeroInt()
	}

	// Total Eden given
	totalEdenGiven := sdk.ZeroInt()
	totalDexRewardsGiven := sdk.ZeroInt()

	// Iterate all delegated validators
	k.stk.IterateDelegations(ctx, delAdr, func(index int64, del stypes.DelegationI) (stop bool) {
		valAddr := del.GetValidatorAddr()
		// Get validator
		val := k.stk.Validator(ctx, valAddr)
		// Get commission rate
		comm_rate := val.GetCommission()
		// Get delegator share
		shares := del.GetShares()
		// Get token amount delegated
		delegatedAmt := val.TokensFromSharesTruncated(shares)

		//-----------------------------
		// Eden commission
		//-----------------------------
		// to give = delegated amount / total delegation * newly minted eden * commission rate
		edenCommission := delegatedAmt.QuoInt(totalDelegationAmt).MulInt(newUncommittedAmt).Mul(comm_rate)
		// Sum total commission given
		totalEdenGiven = totalEdenGiven.Add(edenCommission.TruncateInt())
		//-----------------------------

		//-----------------------------
		// Dex rewards commission
		//-----------------------------
		// to give = delegated amount / total delegation * newly minted eden * commission rate
		dexRewardsCommission := delegatedAmt.QuoInt(totalDelegationAmt).Mul(dexRewards).Mul(comm_rate)
		// Sum total commission given
		totalDexRewardsGiven = totalDexRewardsGiven.Add(dexRewardsCommission.TruncateInt())
		//-----------------------------

		// increase uncomitted token amount of validator's commitment
		k.UpdateTokensForValidator(ctx, valAddr.String(), edenCommission.TruncateInt(), dexRewardsCommission)

		return false
	})

	return totalEdenGiven, totalDexRewardsGiven
}

// withdraw rewards
// Eden, EdenBoost and Elys to USDC
func (k Keeper) ProcessWithdrawRewards(ctx sdk.Context, delegator string) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return err
	}

	// Get commitments
	commitments, bfound := k.cmk.GetCommitments(ctx, delegator)
	if !bfound {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to find commitment")
	}

	// Eden
	uncommittedEden, _ := commitments.GetUncommittedTokensForDenom(ptypes.Eden)
	// Eden B
	uncommittedEdenB, _ := commitments.GetUncommittedTokensForDenom(ptypes.EdenB)
	// Elys
	uncommittedElys, _ := commitments.GetUncommittedTokensForDenom(ptypes.Elys)

	// Withdraw Eden
	err = k.cmk.ProcessWithdrawTokens(ctx, delegator, ptypes.Eden, uncommittedEden.Amount)

	// Withdraw Eden Boost
	err = k.cmk.ProcessWithdrawTokens(ctx, delegator, ptypes.EdenB, uncommittedEdenB.Amount)

	// Convert Elys to USDC
	// Conversion is done inside commmitment module, later on we will have swap module
	err = k.cmk.ProcessWithdrawElysTokens(ctx, delegator, ptypes.Elys, uncommittedElys.Amount)

	return err
}

// withdraw validator commission
// Eden, EdenBoost and Elys to USDC
func (k Keeper) ProcessWithdrawValidatorCommission(ctx sdk.Context, delegator string, validator string) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return err
	}

	// Check validator address
	_, err = sdk.ValAddressFromBech32(validator)
	if err != nil {
		return err
	}

	// Get commitments
	commitments, bfound := k.cmk.GetCommitments(ctx, validator)
	if !bfound {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to find commitment")
	}

	// Eden
	uncommittedEden, _ := commitments.GetUncommittedTokensForDenom(ptypes.Eden)
	// Eden B
	uncommittedEdenB, _ := commitments.GetUncommittedTokensForDenom(ptypes.EdenB)
	// Elys
	uncommittedElys, _ := commitments.GetUncommittedTokensForDenom(ptypes.Elys)

	// Withdraw Eden
	err = k.cmk.ProcessWithdrawValidatorCommission(ctx, delegator, validator, ptypes.Eden, uncommittedEden.Amount)

	// Withdraw Eden Boost
	err = k.cmk.ProcessWithdrawValidatorCommission(ctx, delegator, validator, ptypes.EdenB, uncommittedEdenB.Amount)

	// Convert Elys to USDC
	// Conversion is done inside commmitment module, later on we will have swap module
	err = k.cmk.ProcessWithdrawValidatorElysCommission(ctx, delegator, validator, ptypes.Elys, uncommittedElys.Amount)

	return err
}
