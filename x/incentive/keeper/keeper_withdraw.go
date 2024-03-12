package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// Increase unclaimed token amount for the corresponding validator
func (k Keeper) UpdateTokensForValidator(ctx sdk.Context, validator string, newUnclaimedEdenTokens math.Int, dexRewards sdk.Dec, baseCurrency string) {
	commitments := k.cmk.GetCommitments(ctx, validator)

	// Update Eden amount
	commitments.AddRewardsUnclaimed(sdk.NewCoin(ptypes.Eden, newUnclaimedEdenTokens))

	// Update USDC amount
	commitments.AddRewardsUnclaimed(sdk.NewCoin(baseCurrency, dexRewards.TruncateInt()))

	// Update commmitment
	k.cmk.SetCommitments(ctx, commitments)
}

// Give commissions to validators
func (k Keeper) GiveCommissionToValidators(ctx sdk.Context, delegator string, totalDelegationAmt math.Int, newUnclaimedAmt math.Int, dexRewards sdk.Dec, baseCurrency string) (math.Int, math.Int) {
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
		if !val.IsBonded() {
			return
		}

		// Get commission rate
		commRate := val.GetCommission()
		// Get delegator share
		shares := del.GetShares()
		// Get token amount delegated
		delAmount := val.TokensFromSharesTruncated(shares)

		//-----------------------------
		// Eden commission
		//-----------------------------
		// to give = delegated amount / total delegation * newly minted eden * commission rate
		edenCommission := delAmount.QuoInt(totalDelegationAmt).MulInt(newUnclaimedAmt).Mul(commRate)

		// Sum total commission given
		totalEdenGiven = totalEdenGiven.Add(edenCommission.TruncateInt())
		//-----------------------------

		//-----------------------------
		// Dex rewards commission
		//-----------------------------
		// to give = delegated amount / total delegation * newly minted eden * commission rate
		dexRewardsCommission := delAmount.QuoInt(totalDelegationAmt).Mul(dexRewards).Mul(commRate)
		// Sum total commission given
		totalDexRewardsGiven = totalDexRewardsGiven.Add(dexRewardsCommission.TruncateInt())
		//-----------------------------

		// increase uncomitted token amount of validator's commitment
		k.UpdateTokensForValidator(ctx, valAddr.String(), edenCommission.TruncateInt(), dexRewardsCommission, baseCurrency)

		return false
	})

	return totalEdenGiven, totalDexRewardsGiven
}

// Deduct rewards per program per denom
func (k Keeper) CalcAmountSubbucketsPerProgram(ctx sdk.Context, delegator string, denom string, withdrawType commitmenttypes.EarnType, commitments commitmenttypes.Commitments) math.Int {
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
	case commitmenttypes.EarnType_LP_MINING_PROGRAM:
		unclaimed = commitments.GetLPMiningSubBucketRewardUnclaimedForDenom(denom)
	case commitmenttypes.EarnType_ALL_PROGRAM:
		unclaimed = commitments.GetRewardUnclaimedForDenom(denom)
	}

	return unclaimed
}

// withdraw rewards
// Eden, EdenBoost and Elys to USDC
func (k Keeper) ProcessWithdrawRewards(ctx sdk.Context, delegator string, withdrawType commitmenttypes.EarnType) error {
	_, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return err
	}

	// Get commitments
	commitments := k.cmk.GetCommitments(ctx, delegator)

	// Claim Eden
	// ---------------------------------------------------
	unclaimed := k.CalcAmountSubbucketsPerProgram(ctx, delegator, ptypes.Eden, withdrawType, commitments)
	if !unclaimed.IsZero() {
		err = k.cmk.RecordClaimReward(ctx, delegator, ptypes.Eden, unclaimed, withdrawType)
		if err != nil {
			return err
		}
	}

	// Claim EdenB
	// ---------------------------------------------------
	unclaimed = k.CalcAmountSubbucketsPerProgram(ctx, delegator, ptypes.EdenB, withdrawType, commitments)
	if !unclaimed.IsZero() {
		err = k.cmk.RecordClaimReward(ctx, delegator, ptypes.EdenB, unclaimed, withdrawType)
		if err != nil {
			return err
		}
	}

	// Claim USDC
	// ---------------------------------------------------
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Get available usdc amount can be withdraw
	unclaimedUsdc := k.CalcAmountSubbucketsPerProgram(ctx, delegator, baseCurrency, withdrawType, commitments)
	if unclaimedUsdc.IsZero() {
		return nil
	}
	// Get dex revenue wallet
	revenueCollector := k.authKeeper.GetModuleAccount(ctx, k.dexRevCollectorName)

	// Revenue wallet usdc balance
	usdcBalance := k.bankKeeper.GetBalance(ctx, revenueCollector.GetAddress(), baseCurrency)

	// Balance check
	if unclaimedUsdc.GT(usdcBalance.Amount) {
		return errorsmod.Wrapf(types.ErrIntOverflowTx, "Amount excceed: %d", unclaimedUsdc)
	}

	// All dex rewards are only paid in USDC
	// This function call will deduct the accounting in commitment module only.
	err = k.cmk.RecordClaimReward(ctx, delegator, ptypes.BaseCurrency, unclaimedUsdc, withdrawType)
	if err != nil {
		return errorsmod.Wrapf(types.ErrIntOverflowTx, "Internal error with amount: %d", unclaimedUsdc)
	}

	// Get Bech32 address for creator
	addr, err := sdk.AccAddressFromBech32(commitments.Creator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// Set withdraw usdc amount
	revenue := sdk.NewCoin(baseCurrency, unclaimedUsdc)
	// Transfer revenue from a single wallet of DEX revenue wallet to user.
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.dexRevCollectorName, addr, sdk.NewCoins(revenue))
}

// Update commitments for validator commission
// Eden, EdenBoost and USDC
func (k Keeper) RecordWithdrawValidatorCommission(ctx sdk.Context, delegator string, validator string) error {
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
	unclaimed := commitments.GetRewardUnclaimedForDenom(ptypes.Eden)
	if !unclaimed.IsZero() {
		err = k.cmk.RecordWithdrawValidatorCommission(ctx, delegator, validator, ptypes.Eden, unclaimed)
		if err != nil {
			return err
		}
	}

	// EdenB
	unclaimed = commitments.GetRewardUnclaimedForDenom(ptypes.EdenB)
	if !unclaimed.IsZero() {
		err = k.cmk.RecordWithdrawValidatorCommission(ctx, delegator, validator, ptypes.EdenB, unclaimed)
		if err != nil {
			return err
		}
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// USDC
	unclaimedUsdc := commitments.GetRewardUnclaimedForDenom(baseCurrency)
	if unclaimedUsdc.IsZero() {
		return nil
	}
	// Get dex revenue wallet
	revenueCollector := k.authKeeper.GetModuleAccount(ctx, k.dexRevCollectorName)

	// Revenue wallet usdc balance
	usdcBalance := k.bankKeeper.GetBalance(ctx, revenueCollector.GetAddress(), baseCurrency)

	// Balance check
	if unclaimedUsdc.GT(usdcBalance.Amount) {
		return errorsmod.Wrapf(types.ErrIntOverflowTx, "Amount excceed: %d", unclaimedUsdc)
	}

	// This function call will deduct the accounting in commitment module only.
	err = k.cmk.RecordClaimReward(ctx, validator, baseCurrency, unclaimedUsdc, commitmenttypes.EarnType_ALL_PROGRAM)
	if err != nil {
		return err
	}

	// Get Bech32 address for delegator
	addr, err := sdk.AccAddressFromBech32(delegator)
	if err != nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidAddress, "unable to convert address from bech32")
	}

	// Set withdraw usdc amount
	revenue := sdk.NewCoin(baseCurrency, unclaimedUsdc)
	// Transfer revenue from a single wallet of DEX revenue wallet to user's wallet.
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, k.dexRevCollectorName, addr, sdk.NewCoins(revenue))
}
