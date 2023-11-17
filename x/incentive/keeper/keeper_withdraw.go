package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Increase unclaimed token amount for the corresponding validator
func (k Keeper) UpdateTokensForValidator(ctx sdk.Context, validator string, newUnclaimedEdenTokens sdk.Int, dexRewards sdk.Dec, baseCurrency string) {
	commitments := k.cmk.GetCommitments(ctx, validator)

	// Update Eden amount
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))

	// Update USDC amount
	commitments.AddRewardsUnclaimed(sdk.NewCoin(baseCurrency, dexRewards.TruncateInt()))

	// Update commmitment
	k.cmk.SetCommitments(ctx, commitments)
}

// Give commissions to validators
func (k Keeper) GiveCommissionToValidators(ctx sdk.Context, delegator string, totalDelegationAmt sdk.Int, newUnclaimedAmt sdk.Int, dexRewards sdk.Dec, baseCurrency string) (sdk.Int, sdk.Int) {
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
		k.UpdateTokensForValidator(ctx, valAddr.String(), edenCommission.TruncateInt(), dexRewardsCommission, baseCurrency)

		return false
	})

	return totalEdenGiven, totalDexRewardsGiven
}

// withdraw rewards
// Eden, EdenBoost and Elys to USDC
func (k Keeper) ProcessWithdrawRewards(ctx sdk.Context, delegator string, denom string, withdrawType commitmenttypes.EarnType) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return err
	}

	// Get commitments
	commitments := k.cmk.GetCommitments(ctx, delegator)

	// Eden
	if denom == ptypes.Eden || denom == ptypes.EdenB {
		unclaimed := sdk.ZeroInt()
		switch withdrawType {
		case commitmenttypes.EarnType_ELYS_PROGRAM:
			unclaimed = commitments.GetElysSubBucketRewardUnclaimedForDenom(denom)
		case commitmenttypes.EarnType_EDEN_PROGRAM:
			unclaimed = commitments.GetEdenSubBucketRewardUnclaimedForDenom(denom)
		case commitmenttypes.EarnType_EDENB_PROGRAM:
			unclaimed = commitments.GetEdenBSubBucketRewardUnclaimedForDenom(denom)
		case commitmenttypes.EarnType_USDC_PROGRAM:
			unclaimed = commitments.GetUsdcSubBucketRewardUnclaimedForDenom(denom)
		default:
			unclaimed = commitments.GetRewardUnclaimedForDenom(denom)
		}

		if unclaimed.IsZero() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
		}
		// Claim Eden pr Eden boost from Unclaimed state
		return k.cmk.RecordClaimReward(ctx, delegator, denom, unclaimed, withdrawType)
	}

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// USDC
	unclaimedUsdc := commitments.GetRewardUnclaimedForDenom(baseCurrency)
	if unclaimedUsdc.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}
	// Get dex revenue wallet
	revenueCollector := k.authKeeper.GetModuleAccount(ctx, k.dexRevCollectorName)

	// Revenue wallet usdc balance
	usdcBalance := k.bankKeeper.GetBalance(ctx, revenueCollector.GetAddress(), baseCurrency)

	// Balance check
	if unclaimedUsdc.GT(usdcBalance.Amount) {
		return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Amount excceed: %d", unclaimedUsdc)
	}

	// All dex rewards are only paid in USDC
	// TODO:
	// USDC denom is still dummy until we have real USDC in our chain.
	// This function call will deduct the accounting in commitment module only.
	err = k.cmk.RecordWithdrawUSDC(ctx, delegator, baseCurrency, unclaimedUsdc)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Internal error with amount: %d", unclaimedUsdc)
	}

	// Get Bech32 address for creator
	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// Set withdraw usdc amount
	// TODO
	// USDC denom is still dummy
	revenue := sdk.NewCoin(baseCurrency, unclaimedUsdc)
	// Transfer revenue from a single wallet of DEX revenue wallet to user.
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.dexRevCollectorName, addr, sdk.NewCoins(revenue))
	if err != nil {
		panic(err)
	}

	return err
}

// Update commitments for validator commission
// Eden, EdenBoost and USDC
func (k Keeper) RecordWithdrawValidatorCommission(ctx sdk.Context, delegator string, validator string, denom string) error {
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
	commitments := k.cmk.GetCommitments(ctx, validator)

	// Eden
	if denom == ptypes.Eden || denom == ptypes.EdenB {
		unclaimed := commitments.GetRewardUnclaimedForDenom(denom)
		if unclaimed.IsZero() {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
		}
		// Withdraw Eden or EdenB
		return k.cmk.RecordWithdrawValidatorCommission(ctx, delegator, validator, denom, unclaimed)
	}

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// USDC
	unclaimedUsdc := commitments.GetRewardUnclaimedForDenom(baseCurrency)
	if unclaimedUsdc.IsZero() {

		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "balance not available")
	}
	// Get dex revenue wallet
	revenueCollector := k.authKeeper.GetModuleAccount(ctx, k.dexRevCollectorName)

	// Revenue wallet usdc balance
	usdcBalance := k.bankKeeper.GetBalance(ctx, revenueCollector.GetAddress(), baseCurrency)

	// Balance check
	if unclaimedUsdc.GT(usdcBalance.Amount) {
		return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Amount excceed: %d", unclaimedUsdc)
	}

	// TODO:
	// USDC denom is still dummy until we have real USDC in our chain.
	// This function call will deduct the accounting in commitment module only.
	err = k.cmk.RecordWithdrawUSDC(ctx, validator, baseCurrency, unclaimedUsdc)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrIntOverflowTx, "Internal error with amount: %d", unclaimedUsdc)
	}

	// Get Bech32 address for delegator
	addr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// Set withdraw usdc amount
	// TODO
	// USDC denom is still dummy
	revenue := sdk.NewCoin(baseCurrency, unclaimedUsdc)
	// Transfer revenue from a single wallet of DEX revenue wallet to user's wallet.
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.dexRevCollectorName, addr, sdk.NewCoins(revenue))
	if err != nil {
		panic(err)
	}

	return err
}
