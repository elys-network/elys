package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) AddCollateralToMtp(ctx sdk.Context, msg *types.MsgAddCollateral) (*types.MsgAddCollateralResponse, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Load existing position
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, creator, msg.Id)
	if err != nil {
		return nil, err
	}
	senderAddress, err := sdk.AccAddressFromBech32(mtp.Address)
	if err != nil {
		return nil, err
	}

	// Fetch the pool associated with the given pool ID.
	pool, found := k.OpenDefineAssetsChecker.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, mtp.TradingAsset)
	}

	// Check if the pool is enabled.
	if !k.OpenDefineAssetsChecker.IsPoolEnabled(ctx, mtp.AmmPoolId) {
		return nil, errorsmod.Wrap(types.ErrMTPDisabled, mtp.TradingAsset)
	}

	// Fetch the corresponding AMM (Automated Market Maker) pool.
	ammPool, err := k.OpenDefineAssetsChecker.GetAmmPool(ctx, mtp.AmmPoolId, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	collateralCoin := sdk.NewCoin(mtp.CollateralAsset, msg.Amount)

	if !k.bankKeeper.HasBalance(ctx, senderAddress, collateralCoin) {
		return nil, types.ErrBalanceNotAvailable
	}

	//collateralAmountDec := sdk.NewDecFromBigInt(msg.Amount.BigInt())
	liabilitiesDec := collateralCoin.Amount.ToLegacyDec()
	// As we are adding to current position, liabilities will decrease

	// If collateral asset is not base currency, should calculate liability in base currency with the given out.
	// Liability has to be in base currency
	if mtp.CollateralAsset != baseCurrency {
		// ATOM amount
		etaAmt := liabilitiesDec.TruncateInt()
		etaAmtToken := sdk.NewCoin(mtp.CollateralAsset, etaAmt)
		// Calculate base currency amount given atom out amount and we use it liabilty amount in base currency
		liabilityAmt, err := k.OpenDefineAssetsChecker.EstimateSwapGivenOut(ctx, etaAmtToken, baseCurrency, ammPool)
		if err != nil {
			return nil, err
		}

		liabilitiesDec = sdk.NewDecFromInt(liabilityAmt)
	}

	// If position is short, liabilities should be swapped to liabilities asset
	if mtp.Position == types.Position_SHORT {
		liabilitiesAmtTokenIn := sdk.NewCoin(baseCurrency, liabilitiesDec.TruncateInt())
		liabilitiesAmt, err := k.OpenDefineAssetsChecker.EstimateSwap(ctx, liabilitiesAmtTokenIn, mtp.LiabilitiesAsset, ammPool)
		if err != nil {
			return nil, err
		}

		liabilitiesDec = sdk.NewDecFromInt(liabilitiesAmt)
	}

	// Add collateral, decrease liability
	mtp.Collateral = mtp.Collateral.Add(collateralCoin.Amount)

	mtp.Liabilities = mtp.Liabilities.Sub(sdk.NewIntFromBigInt(liabilitiesDec.TruncateInt().BigInt()))

	// calculate mtp take profit custody, delta y_tp_c = delta x_l / take profit price (take profit custody = liabilities / take profit price)
	mtp.TakeProfitCustody = types.CalcMTPTakeProfitCustody(&mtp)

	// calculate mtp take profit liabilities, delta x_tp_l = delta y_tp_c * current price (take profit liabilities = take profit custody * current price)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, &mtp, baseCurrency)
	if err != nil {
		return nil, err
	}

	h, err := k.GetMTPHealth(ctx, mtp, ammPool, baseCurrency) // set mtp in func or return h?
	if err != nil {
		return nil, err
	}
	mtp.MtpHealth = h

	ammPoolAddr, err := sdk.AccAddressFromBech32(ammPool.Address)
	if err != nil {
		return nil, err
	}

	collateralCoins := sdk.NewCoins(collateralCoin)
	err = k.bankKeeper.SendCoins(ctx, senderAddress, ammPoolAddr, collateralCoins)

	if err != nil {
		return nil, err
	}

	err = pool.UpdateBalance(ctx, mtp.CollateralAsset, collateralCoin.Amount, true, mtp.Position)
	if err != nil {
		return nil, err
	}

	// All liability has to be in liabilities asset
	// Decrease pool liability
	err = pool.UpdateLiabilities(ctx, mtp.LiabilitiesAsset, liabilitiesDec.TruncateInt(), false, mtp.Position)
	if err != nil {
		return nil, err
	}

	// All take profit liability has to be in liabilities asset
	err = pool.UpdateTakeProfitLiabilities(ctx, mtp.LiabilitiesAsset, mtp.TakeProfitLiabilities, true, mtp.Position)
	if err != nil {
		return nil, err
	}

	// All take profit custody has to be in custody asset
	err = pool.UpdateTakeProfitCustody(ctx, mtp.CustodyAsset, mtp.TakeProfitCustody, true, mtp.Position)
	if err != nil {
		return nil, err
	}

	k.SetPool(ctx, pool)

	// calc and update open price
	err = k.OpenChecker.UpdateOpenPrice(ctx, &mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Update the pool health.
	if err = k.OpenDefineAssetsChecker.UpdatePoolHealth(ctx, &pool); err != nil {
		return nil, err
	}

	// Update the MTP health.
	lr, err := k.OpenDefineAssetsChecker.GetMTPHealth(ctx, mtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.OpenDefineAssetsChecker.GetSafetyFactor(ctx)
	if lr.LTE(safetyFactor) {
		return nil, types.ErrMTPUnhealthy
	}

	// Update consolidated collateral amount
	err = k.OpenDefineAssetsChecker.CalcMTPConsolidateCollateral(ctx, &mtp, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Calculate consolidate liability and update consolidate leverage
	mtp.ConsolidateLeverage = types.CalcMTPConsolidateLiability(&mtp)

	// Set MTP
	err = k.OpenDefineAssetsChecker.SetMTP(ctx, &mtp)
	if err != nil {
		return nil, err
	}

	k.EmitOpenEvent(ctx, &mtp)

	if k.hooks != nil {
		k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator)
	}

	return &types.MsgAddCollateralResponse{}, nil
}
