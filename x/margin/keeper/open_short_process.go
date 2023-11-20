package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) ProcessOpenShort(ctx sdk.Context, mtp *types.MTP, leverage sdk.Dec, eta sdk.Dec, collateralAmountDec sdk.Dec, poolId uint64, msg *types.MsgOpen, baseCurrency string) (*types.MTP, error) {
	// Determine the trading asset.
	tradingAsset := types.GetTradingAsset(msg.CollateralAsset, msg.BorrowAsset, baseCurrency)

	// Fetch the pool associated with the given pool ID.
	pool, found := k.OpenShortChecker.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, tradingAsset)
	}

	// Check if the pool is enabled.
	if !k.OpenShortChecker.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, tradingAsset)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.OpenShortChecker.GetAmmPool(ctx, poolId, tradingAsset)
	if err != nil {
		return nil, err
	}

	// Calculate the leveraged amount based on the collateral provided and the leverage.
	leveragedAmount := sdk.NewInt(collateralAmountDec.Mul(leverage).TruncateInt().Int64())

	if msg.CollateralAsset != baseCurrency {
		return nil, sdkerrors.Wrap(types.ErrInvalidBorrowingAsset, "collateral must be base currency")
	}

	custodyAmtToken := sdk.NewCoin(baseCurrency, leveragedAmount)
	borrowingAmount, err := k.OpenShortChecker.EstimateSwapGivenOut(ctx, custodyAmtToken, msg.BorrowAsset, ammPool)
	if err != nil {
		return nil, err
	}

	// check the balance
	if !types.HasSufficientPoolBalance(ammPool, baseCurrency, borrowingAmount) {
		return nil, sdkerrors.Wrap(types.ErrBorrowTooHigh, borrowingAmount.String())
	}

	// Check minimum liabilities.
	collateralTokenAmt := sdk.NewCoin(msg.CollateralAsset, msg.CollateralAmount)
	err = k.OpenShortChecker.CheckMinLiabilities(ctx, collateralTokenAmt, eta, pool, ammPool, msg.BorrowAsset)
	if err != nil {
		return nil, err
	}

	// Calculate custody amount.
	leveragedAmtTokenIn := sdk.NewCoin(msg.BorrowAsset, borrowingAmount)
	custodyAmount, err := k.OpenShortChecker.EstimateSwap(ctx, leveragedAmtTokenIn, baseCurrency, ammPool)
	if err != nil {
		return nil, err
	}

	// Ensure the AMM pool has enough balance.
	if !types.HasSufficientPoolBalance(ammPool, baseCurrency, custodyAmount) {
		return nil, sdkerrors.Wrap(types.ErrCustodyTooHigh, custodyAmount.String())
	}

	// if position is short then override the custody asset to the base currency
	if mtp.Position == types.Position_SHORT {
		mtp.Custodies = []sdk.Coin{sdk.NewCoin(baseCurrency, sdk.NewInt(0))}
		mtp.BorrowInterestPaidCustodies = []sdk.Coin{sdk.NewCoin(baseCurrency, sdk.NewInt(0))}
		mtp.TakeProfitCustodies = []sdk.Coin{sdk.NewCoin(baseCurrency, sdk.NewInt(0))}
	}

	// Borrow the asset the user wants to short.
	err = k.OpenShortChecker.Borrow(ctx, msg.CollateralAsset, baseCurrency, msg.CollateralAmount, custodyAmount, mtp, &ammPool, &pool, eta, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.OpenShortChecker.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Take custody from the pool balance.
	if err = k.OpenShortChecker.TakeInCustody(ctx, *mtp, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.OpenShortChecker.UpdateMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.OpenShortChecker.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Update consolidated collateral amount
	k.OpenShortChecker.CalcMTPConsolidateCollateral(ctx, mtp, baseCurrency)

	// Calculate consolidate liabiltiy and update consolidate leverage
	mtp.ConsolidateLeverage = types.CalcMTPConsolidateLiability(mtp)

	// Set MTP
	k.OpenShortChecker.SetMTP(ctx, mtp)

	// Return the updated Margin Trading Position (MTP).
	return mtp, nil
}
