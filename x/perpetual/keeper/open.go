package keeper

import (
	"errors"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen, skipTriggerCheck bool) (*types.MsgOpenResponse, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("PoolId: %d", msg.PoolId))
	}

	tradingAsset, err := pool.GetTradingAsset(baseCurrency)
	if err != nil {
		return nil, err
	}

	if err := msg.ValidatePosition(tradingAsset, baseCurrency); err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)
	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, tradingAsset)
	if err != nil {
		return nil, err
	}
	if tradingAssetPrice.IsZero() {
		return nil, errors.New("trading asset price is zero while opening perpetual")
	}

	if err = msg.ValidateTakeProfitAndStopLossPrice(params, tradingAssetPrice); err != nil {
		return nil, err
	}

	if err = k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	// check if existing mtp to consolidate
	existingMtp := k.GetExistingPosition(ctx, msg)

	if existingMtp == nil {
		// opening new position
		if msg.Leverage.LTE(math.LegacyOneDec()) {
			return nil, errors.New("cannot open new position with leverage <= 1")
		}
		// Check if max positions are exceeded as we are opening new position, not updating old position
		if err = k.CheckMaxOpenPositions(ctx, msg.PoolId); err != nil {
			return nil, err
		}
	} else if msg.Leverage.Equal(math.LegacyZeroDec()) {
		// adding collateral to existing position (when leverage > 1, we leave the case for modifying old position)
		// Enforce collateral addition (for leverage 1) without modifying take profit price
		msg.TakeProfitPrice = existingMtp.TakeProfitPrice
	}

	if err = k.CheckLowPoolHealthAndMinimumCustody(ctx, msg.PoolId, msg.Position); err != nil {
		return nil, err
	}

	proxyLeverage := GetProxyLeverage(msg.Position, msg.Leverage, pool)

	// Define the assets
	custodyAsset, liabilitiesAsset, err := msg.GetCustodyAndLiabilitiesAsset(tradingAsset, baseCurrency)
	if err != nil {
		return nil, err
	}

	ammPool, err := k.GetAmmPool(ctx, msg.PoolId)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "amm pool id %d", msg.PoolId)
	}

	// Initialize a new Perpetual Trading Position (MTP).
	mtp := types.NewMTP(ctx, msg.Creator, msg.Collateral.Denom, tradingAsset, liabilitiesAsset, custodyAsset, msg.Position, msg.TakeProfitPrice, msg.PoolId)

	totalPerpFeesCoins, err := k.ProcessOpen(ctx, &pool, &ammPool, mtp, proxyLeverage, msg.PoolId, msg, baseCurrency)
	if err != nil {
		return nil, err
	}

	if existingMtp != nil {
		return k.OpenConsolidate(ctx, existingMtp, mtp, msg, tradingAsset, totalPerpFeesCoins, skipTriggerCheck)
	}

	if err = k.CheckLowPoolHealthAndMinimumCustody(ctx, msg.PoolId, msg.Position); err != nil {
		return nil, err
	}

	// should not be checked before OpenConsolidate
	denomPrice, err := k.GetDenomPrice(ctx, tradingAsset)
	if err != nil {
		return nil, err
	}
	if mtp.GetMTPValue(denomPrice).LT(params.MinimumNotionalValue) {
		return nil, fmt.Errorf("not enough notional value for the mtp: minimum %s, mtp: %s", params.MinimumNotionalValue.String(), mtp.GetMTPValue(denomPrice).String())
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		err = k.hooks.AfterPerpetualPositionOpen(ctx, ammPool, pool, creator)
		if err != nil {
			return nil, err
		}
	}

	perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd := k.GetPerpFeesInUSD(ctx, totalPerpFeesCoins)

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventOpen,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("owner", mtp.Address),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("collateral", mtp.Collateral.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("custody", mtp.Custody.String()),
		sdk.NewAttribute("mtp_health", mtp.MtpHealth.String()),
		sdk.NewAttribute("stop_loss_price", mtp.StopLossPrice.String()),
		sdk.NewAttribute("take_profit_price", mtp.TakeProfitPrice.String()),
		sdk.NewAttribute("take_profit_borrow_factor", mtp.TakeProfitBorrowFactor.String()),
		sdk.NewAttribute("funding_fee_paid_custody", mtp.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_custody", mtp.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute("open_price", mtp.OpenPrice.String()),
		sdk.NewAttribute(types.AttributeKeyPerpFee, perpFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeySlippage, slippageFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeyWeightBreakingFee, weightBreakingFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeTakerFees, takerFeesInUsd.String()),
	))

	return &types.MsgOpenResponse{
		Id: mtp.Id,
	}, nil
}

func GetProxyLeverage(position types.Position, leverage math.LegacyDec, pool types.Pool) math.LegacyDec {
	// Determine the maximum leverage available for this pool and compute the effective leverage to be used.
	// values for leverage other than 0 or  >1 are invalidated in validate basic
	proxyLeverage := math.LegacyMinDec(leverage, pool.LeverageMax)

	// just adding collateral
	if leverage.IsZero() {
		proxyLeverage = math.LegacyOneDec()
	} else {
		// opening position, for Short we add 1 because, say atom price 5 usdc, collateral 100 usdc, leverage 5, then liabilities will be 80 atom worth 400 usdc which would be position size
		// User would be expecting position size of 100 atom / 500 usdc. So we increase the leverage from 5 to 6
		// Because of this effective leverage for short has to be reduced by 1 in query
		if position == types.Position_SHORT {
			proxyLeverage = proxyLeverage.Add(math.LegacyOneDec())
		}
		// We don't need to do this for LONG as it gives desired position
	}

	return proxyLeverage
}
