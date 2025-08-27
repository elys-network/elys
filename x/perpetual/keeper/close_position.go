package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k Keeper) ClosePosition(ctx sdk.Context, msg *types.MsgClose) (types.MTP, math.Int, math.LegacyDec, math.Int, math.Int, math.Int, math.Int, math.Int, bool, bool, types.PerpetualFees, math.LegacyDec, sdk.Coin, math.Int, math.Int, error) {
	// Retrieve MTP
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	zeroPerpFees := types.NewPerpetualFeesWithEmptyCoins()
	mtp, err := k.GetMTP(ctx, msg.PoolId, creator, msg.Id)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), sdk.Coin{}, math.ZeroInt(), math.ZeroInt(), err
	}

	userClosingRatio := msg.Amount.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec())
	if mtp.Position == types.Position_SHORT {
		userClosingRatio = msg.Amount.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
	}
	if userClosingRatio.GT(math.LegacyOneDec()) {
		userClosingRatio = math.LegacyOneDec()
	}

	initialCollateral := sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral)
	initialCustody := mtp.Custody
	initialLiabilities := mtp.Liabilities

	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), initialCollateral, initialCustody, initialLiabilities, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", mtp.AmmPoolId))
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), initialCollateral, initialCustody, initialLiabilities, err
	}

	// this also handles edge case where bot is unable to close position in time.
	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, closingRatio, err := k.MTPTriggerChecksAndUpdates(ctx, &mtp, &pool, &ammPool)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), initialCollateral, initialCustody, initialLiabilities, err
	}

	if forceClosed {
		// user didn't want close position fully and it got only partially closed, we return as the value user wanted to close is different now
		// OR it got fully liquidated
		if (!userClosingRatio.Equal(math.LegacyOneDec()) && !closingRatio.Equal(math.LegacyOneDec())) || closingRatio.Equal(math.LegacyOneDec()) {
			return mtp, repayAmt, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, initialCollateral, initialCustody, initialLiabilities, nil
		}
	}

	// Should be reset after MTPTriggerChecksAndUpdates
	closingRatio = userClosingRatio

	// Estimate swap and repay
	repayAmt, returnAmt, perpFees, closingPrice, _, err := k.EstimateAndRepay(ctx, &mtp, &pool, &ammPool, closingRatio, false)
	if err != nil {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), allInterestsPaid, forceClosed, perpetualFeesCoins, math.LegacyZeroDec(), initialCollateral, initialCustody, initialLiabilities, err
	}
	perpetualFeesCoins = perpetualFeesCoins.Add(perpFees)

	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		err = k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, pool, creator, closingRatio, mtp.Id)
		if err != nil {
			return mtp, math.Int{}, math.LegacyDec{}, math.Int{}, math.Int{}, math.Int{}, math.Int{}, math.Int{}, allInterestsPaid, forceClosed, zeroPerpFees, math.LegacyZeroDec(), initialCollateral, initialCustody, initialLiabilities, err
		}
	}

	return mtp, repayAmt, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, initialCollateral, initialCustody, initialLiabilities, nil
}
