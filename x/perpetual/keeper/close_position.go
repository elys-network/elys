package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) ClosePosition(ctx sdk.Context, msg *types.MsgClose) (types.MTP, math.Int, math.LegacyDec, math.Int, math.Int, math.Int, math.Int, math.Int, bool, bool, types.PerpetualFees, math.LegacyDec, sdk.Coin, error) {
	// Retrieve MTP
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	zeroPerpFees := types.NewPerpetualFeesWithEmptyCoins()
	mtp, err := k.GetMTP(ctx, msg.PoolId, creator, msg.Id)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), sdk.Coin{}, err
	}

	closingCollatoral := sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral)

	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), closingCollatoral, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", mtp.AmmPoolId))
	}

	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), closingCollatoral, err
	}

	// this also handles edge case where bot is unable to close position in time.
	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, err := k.MTPTriggerChecksAndUpdates(ctx, &mtp, &pool, &ammPool)
	if err != nil {
		return types.MTP{}, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), false, false, zeroPerpFees, math.LegacyZeroDec(), closingCollatoral, err
	}

	if forceClosed {
		return mtp, repayAmt, math.LegacyOneDec(), returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, closingCollatoral, nil
	}

	// Should be declared after SettleMTPBorrowInterestUnpaidLiability and settling funding
	closingRatio := msg.Amount.ToLegacyDec().Quo(mtp.Custody.ToLegacyDec())
	if mtp.Position == types.Position_SHORT {
		closingRatio = msg.Amount.ToLegacyDec().Quo(mtp.Liabilities.ToLegacyDec())
	}
	if closingRatio.GT(math.LegacyOneDec()) {
		closingRatio = math.LegacyOneDec()
	}

	// Estimate swap and repay
	repayAmt, returnAmt, perpFees, closingPrice, err := k.EstimateAndRepay(ctx, &mtp, &pool, &ammPool, closingRatio)
	if err != nil {
		return mtp, math.ZeroInt(), math.LegacyZeroDec(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), math.ZeroInt(), allInterestsPaid, forceClosed, perpetualFeesCoins, math.LegacyZeroDec(), closingCollatoral, err
	}
	perpetualFeesCoins = perpetualFeesCoins.Add(perpFees)

	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		err = k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, pool, creator, closingRatio, mtp.Id)
		if err != nil {
			return mtp, math.Int{}, math.LegacyDec{}, math.Int{}, math.Int{}, math.Int{}, math.Int{}, math.Int{}, allInterestsPaid, forceClosed, zeroPerpFees, math.LegacyZeroDec(), closingCollatoral, err
		}
	}

	return mtp, repayAmt, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, closingCollatoral, nil
}
