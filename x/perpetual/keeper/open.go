package keeper

import (
	"errors"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) Open(ctx sdk.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// Determine the type of position (long or short) and validate assets accordingly.
	switch msg.Position {
	case types.Position_LONG:
		if err := types.CheckLongAssets(msg.Collateral.Denom, msg.TradingAsset, baseCurrency); err != nil {
			return nil, err
		}
	case types.Position_SHORT:
		if err := types.CheckShortAssets(msg.Collateral.Denom, msg.TradingAsset, baseCurrency); err != nil {
			return nil, err
		}
	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPosition, msg.Position.String())
	}

	params := k.GetParams(ctx)
	tradingAssetPrice, err := k.GetAssetPrice(ctx, msg.TradingAsset)
	if err != nil {
		return nil, err
	}
	if tradingAssetPrice.IsZero() {
		return nil, errors.New("trading asset price is zero while opening perpetual")
	}
	ratio := osmomath.BigDecFromDec(msg.TakeProfitPrice).Quo(tradingAssetPrice)
	if msg.Position == types.Position_LONG {
		if ratio.LT(params.GetBigDecMinimumLongTakeProfitPriceRatio()) || ratio.GT(params.GetBigDecMaximumLongTakeProfitPriceRatio()) {
			return nil, fmt.Errorf("take profit price should be between %s and %s times of current market price for long (current ratio: %s)", params.MinimumLongTakeProfitPriceRatio.String(), params.MaximumLongTakeProfitPriceRatio.String(), ratio.String())
		}
		if !msg.StopLossPrice.IsZero() && osmomath.BigDecFromDec(msg.StopLossPrice).GTE(tradingAssetPrice) {
			return nil, fmt.Errorf("stop loss price cannot be greater than equal to tradingAssetPrice for long (Stop loss: %s, asset price: %s)", msg.StopLossPrice.String(), tradingAssetPrice.String())
		}
		// no need to override msg.TakeProfitPrice as the above ratio check it
	}
	if msg.Position == types.Position_SHORT {
		if ratio.GT(params.GetBigDecMaximumShortTakeProfitPriceRatio()) {
			return nil, fmt.Errorf("take profit price should be less than %s times of current market price for short (current ratio: %s)", params.MaximumShortTakeProfitPriceRatio.String(), ratio.String())
		}
		if !msg.StopLossPrice.IsZero() && osmomath.BigDecFromDec(msg.StopLossPrice).LTE(tradingAssetPrice) {
			return nil, fmt.Errorf("stop loss price cannot be less than equal to tradingAssetPrice for short (Stop loss: %s, asset price: %s)", msg.StopLossPrice.String(), tradingAssetPrice.String())
		}
	}

	if err = k.CheckUserAuthorization(ctx, msg); err != nil {
		return nil, err
	}

	// check if existing mtp to consolidate
	existingMtp := k.CheckSameAssetPosition(ctx, msg)

	if existingMtp == nil {
		// opening new position
		if msg.Leverage.LTE(math.LegacyOneDec()) {
			return nil, errors.New("cannot open new position with leverage <= 1")
		}
		// Check if max positions are exceeded as we are opening new position, not updating old position
		if err = k.CheckMaxOpenPositions(ctx); err != nil {
			return nil, err
		}
	} else if msg.Leverage.Equal(math.LegacyZeroDec()) {
		// adding collateral to existing position (when leverage > 1, we leave the case for modifying old position)
		// Enforce collateral addition (for leverage 1) without modifying take profit price
		msg.TakeProfitPrice = existingMtp.TakeProfitPrice
	}

	poolId := msg.PoolId
	// Get pool id, amm pool, and perpetual pool
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "amm pool not found for pool %d", poolId)
	}

	if !ammPool.PoolParams.UseOracle {
		return nil, types.ErrPoolHasToBeOracle
	}

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	if err = k.CheckLowPoolHealthAndMinimumCustody(ctx, poolId); err != nil {
		return nil, err
	}

	mtp, err := k.OpenDefineAssets(ctx, poolId, msg, baseCurrency)
	if err != nil {
		return nil, err
	}

	// calc and update open price
	err = k.UpdateOpenPrice(ctx, mtp)
	if err != nil {
		return nil, err
	}

	if existingMtp != nil {
		return k.OpenConsolidate(ctx, existingMtp, mtp, msg, baseCurrency)
	}

	if err = k.CheckLowPoolHealthAndMinimumCustody(ctx, poolId); err != nil {
		return nil, err
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		// pool values has been updated
		pool, found = k.GetPool(ctx, poolId)
		if !found {
			return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
		}

		err = k.hooks.AfterPerpetualPositionOpen(ctx, ammPool, pool, creator, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return nil, err
		}
	}

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
	))

	return &types.MsgOpenResponse{
		Id: mtp.Id,
	}, nil
}
