package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/margin/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) ProcessOpenLong(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, eta sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen) (*types.MTP, error) {
	// Determine the trading asset.
	tradingAsset := k.OpenLongChecker.GetTradingAsset(msg.CollateralAsset, msg.BorrowAsset)

	// Fetch the pool associated with the given pool ID.
	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, tradingAsset)
	}

	// Check if the pool is enabled.
	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, tradingAsset)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId, tradingAsset)
	if err != nil {
		return nil, err
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// If collateral is not base currency, calculate the borrowing amount in base currency and check the balance
	if msg.CollateralAsset != baseCurrency {
		custodyAmtToken := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)
		borrowingAmount, err := k.OpenLongChecker.EstimateSwapGivenOut(ctx, custodyAmtToken, baseCurrency, ammPool)
		if err != nil {
			return nil, err
		}
		if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, baseCurrency, borrowingAmount) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, borrowingAmount.String())
		}
	} else {
		if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, msg.CollateralAsset, leveragedAmount) {
			return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, leveragedAmount.String())
		}
	}

	// Check minimum liabilities.
	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	err = k.OpenLongChecker.CheckMinLiabilities(ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset)
	if err != nil {
		return nil, err
	}

	// Calculate custody amount.
	leveragedAmtTokenIn := sdk.NewCoin(msg.CollateralAsset, leveragedAmount)
	custodyAmount, err := k.OpenLongChecker.EstimateSwap(ctx, leveragedAmtTokenIn, msg.BorrowAsset, ammPool)
	if err != nil {
		return nil, err
	}

	// If the collateral asset is not base currency, custody amount equals to leverage amount
	if msg.CollateralAsset != baseCurrency {
		custodyAmount = leveragedAmount
	}

	// Ensure the AMM pool has enough balance.
	if !k.OpenLongChecker.HasSufficientPoolBalance(ctx, ammPool, msg.BorrowAsset, custodyAmount) {
		return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	// Borrow the asset the user wants to long.
	err = k.OpenLongChecker.Borrow(ctx, msg.CollateralAsset, msg.BorrowAsset, msg.CollateralAmount, custodyAmount, mtp, &ammPool, &pool, eta)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.OpenLongChecker.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Take custody from the pool balance.
	if err = k.OpenLongChecker.TakeInCustody(ctx, *mtp, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.OpenLongChecker.UpdateMTPHealth(ctx, *mtp, ammPool)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.OpenLongChecker.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Update consolidated collateral amount
	k.OpenLongChecker.CalcMTPConsolidateCollateral(ctx, mtp)

	// Calculate consolidate liabiltiy
	k.OpenLongChecker.CalcMTPConsolidateLiability(ctx, mtp)

	// Set MTP
	k.OpenLongChecker.SetMTP(ctx, mtp)

	return mtp, nil
}
