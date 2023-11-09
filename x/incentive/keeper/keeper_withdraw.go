package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Increase unclaimed token amount for the corresponding validator
func (k Keeper) UpdateTokensForValidator(ctx sdk.Context, validator string, newUnclaimedEdenTokens sdk.Int, dexRewards sdk.Dec) {
	commitments, bfound := k.cmk.GetCommitments(ctx, validator)
	if !bfound {
		return
	}

	// Update Eden amount
	k.UpdateTokensCommitment(&commitments, newUnclaimedEdenTokens, ptypes.Eden)

	// Update USDC amount
	k.UpdateTokensCommitment(&commitments, dexRewards.TruncateInt(), ptypes.BaseCurrency)

	// Update commmitment
	k.cmk.SetCommitments(ctx, commitments)
}

// Give commissions to validators
func (k Keeper) GiveCommissionToValidators(ctx sdk.Context, delegator string, totalDelegationAmt sdk.Int, newUnclaimedAmt sdk.Int, dexRewards sdk.Dec) (sdk.Int, sdk.Int) {
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
		edenCommission := delegatedAmt.QuoInt(totalDelegationAmt).MulInt(newUnclaimedAmt).Mul(comm_rate)

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
func (k Keeper) ProcessWithdrawRewards(ctx sdk.Context, delegator string, denom string) error {
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
	if denom == ptypes.Eden {
		unclaimedEden, bfound := commitments.GetRewardsUnclaimedForDenom(ptypes.Eden)
		if bfound {
			// Withdraw Eden
			return k.cmk.ProcessWithdrawTokens(ctx, delegator, ptypes.Eden, unclaimedEden.Amount)
		}
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}

	if denom == ptypes.EdenB {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}

	// USDC
	unclaimedUsdc, bfound := commitments.GetRewardsUnclaimedForDenom(ptypes.BaseCurrency)
	if bfound {
		// Get dex revenue wallet
		revenueCollector := k.authKeeper.GetModuleAccount(ctx, k.dexRevCollectorName)

		// Revenue wallet usdc balance
		usdcBalance := k.bankKeeper.GetBalance(ctx, revenueCollector.GetAddress(), ptypes.BaseCurrency)

		// Balance check
		if unclaimedUsdc.Amount.GT(usdcBalance.Amount) {
			return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Amount excceed: %d", unclaimedUsdc.Amount)
		}

		// All dex rewards are only paid in USDC
		// TODO:
		// USDC denom is still dummy until we have real USDC in our chain.
		// This function call will deduct the accounting in commitment module only.
		err = k.cmk.ProcessWithdrawUSDC(ctx, delegator, ptypes.BaseCurrency, unclaimedUsdc.Amount)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Internal error with amount: %d", unclaimedUsdc.Amount)
		}

		// Get Bech32 address for creator
		addr, err := sdk.AccAddressFromBech32(commitments.Creator)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
		}

		// Set withdraw usdc amount
		// TODO
		// USDC denom is still dummy
		revenue := sdk.NewCoin(ptypes.BaseCurrency, unclaimedUsdc.Amount)
		// Transfer revenue from a single wallet of DEX revenue wallet to user.
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.dexRevCollectorName, addr, sdk.NewCoins(revenue))
		if err != nil {
			panic(err)
		}
	} else {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}

	return err
}

// Withdraw validator commission
// Eden, EdenBoost and USDC
func (k Keeper) ProcessWithdrawValidatorCommission(ctx sdk.Context, delegator string, validator string, denom string) error {
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
	if denom == ptypes.Eden {
		unclaimedEden, bfound := commitments.GetRewardsUnclaimedForDenom(ptypes.Eden)
		if bfound {
			// Withdraw Eden
			return k.cmk.ProcessWithdrawValidatorCommission(ctx, delegator, validator, ptypes.Eden, unclaimedEden.Amount)
		}
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}

	if denom == ptypes.EdenB {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}

	// USDC
	unclaimedUsdc, bfound := commitments.GetRewardsUnclaimedForDenom(ptypes.BaseCurrency)
	if bfound {
		// Get dex revenue wallet
		revenueCollector := k.authKeeper.GetModuleAccount(ctx, k.dexRevCollectorName)

		// Revenue wallet usdc balance
		usdcBalance := k.bankKeeper.GetBalance(ctx, revenueCollector.GetAddress(), ptypes.BaseCurrency)

		// Balance check
		if unclaimedUsdc.Amount.GT(usdcBalance.Amount) {
			return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Amount excceed: %d", unclaimedUsdc.Amount)
		}

		// TODO:
		// USDC denom is still dummy until we have real USDC in our chain.
		// This function call will deduct the accounting in commitment module only.
		err = k.cmk.ProcessWithdrawUSDC(ctx, validator, ptypes.BaseCurrency, unclaimedUsdc.Amount)
		if err != nil {
			return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Internal error with amount: %d", unclaimedUsdc.Amount)
		}

		// Get Bech32 address for delegator
		addr, err := sdk.AccAddressFromBech32(delegator)
		if err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
		}

		// Set withdraw usdc amount
		// TODO
		// USDC denom is still dummy
		revenue := sdk.NewCoin(ptypes.BaseCurrency, unclaimedUsdc.Amount)
		// Transfer revenue from a single wallet of DEX revenue wallet to user's wallet.
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.dexRevCollectorName, addr, sdk.NewCoins(revenue))
		if err != nil {
			panic(err)
		}
	} else {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}

	return err
}
