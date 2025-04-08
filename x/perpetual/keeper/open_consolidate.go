package keeper

import (
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, baseCurrency string) (*types.MsgOpenResponse, error) {
	poolId := existingMtp.AmmPoolId
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.MTPTriggerChecksAndUpdates(ctx, existingMtp, &pool, &ammPool)
	if err != nil {
		return nil, err
	}

	if forceClosed {
		tradingAssetPrice, err := k.GetAssetPrice(ctx, msg.TradingAsset)
		if err != nil {
			return nil, err
		}
		k.EmitForceClose(ctx, "open_consolidate", *existingMtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, msg.Creator, allInterestsPaid, tradingAssetPrice)
		return &types.MsgOpenResponse{
			Id: existingMtp.Id,
		}, nil
	}

	existingMtp, err = k.OpenConsolidateMergeMtp(ctx, existingMtp, newMtp)
	if err != nil {
		return nil, err
	}

	if !newMtp.Liabilities.IsZero() {
		consolidatedOpenPrice := (existingMtp.Custody.ToLegacyDec().Mul(existingMtp.OpenPrice).Add(newMtp.Custody.ToLegacyDec().Mul(newMtp.OpenPrice))).Quo(existingMtp.Custody.ToLegacyDec().Add(newMtp.Custody.ToLegacyDec()))
		existingMtp.OpenPrice = consolidatedOpenPrice

		consolidatedTakeProfitPrice := existingMtp.Custody.ToLegacyDec().Mul(existingMtp.TakeProfitPrice).Add(newMtp.Custody.ToLegacyDec().Mul(newMtp.TakeProfitPrice)).Quo(existingMtp.Custody.ToLegacyDec().Add(newMtp.Custody.ToLegacyDec()))
		existingMtp.TakeProfitPrice = consolidatedTakeProfitPrice
	}

	existingMtp.TakeProfitCustody = existingMtp.TakeProfitCustody.Add(newMtp.TakeProfitCustody)
	existingMtp.TakeProfitLiabilities = existingMtp.TakeProfitLiabilities.Add(newMtp.TakeProfitLiabilities)

	// no need to update pool's TakeProfitCustody, TakeProfitLiabilities, Custody and Liabilities as it was already in OpenDefineAssets

	existingMtp.MtpHealth, err = k.GetMTPHealth(ctx, *existingMtp, ammPool, baseCurrency)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	safetyFactor := k.GetSafetyFactor(ctx)
	if existingMtp.MtpHealth.LTE(safetyFactor) {
		return nil, errorsmod.Wrapf(types.ErrMTPUnhealthy, "(MtpHealth: %s)", existingMtp.MtpHealth.String())
	}

	stopLossPrice := msg.StopLossPrice
	if msg.StopLossPrice.IsNil() || msg.StopLossPrice.IsZero() {
		stopLossPrice = k.GetLiquidationPrice(ctx, *existingMtp)
	}
	existingMtp.StopLossPrice = stopLossPrice

	// Set existing MTP
	if err = k.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		params := k.GetParams(ctx)
		// The pool value above was sent in pointer so its updated
		err = k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return nil, err
		}
	}

	if err = k.CheckLowPoolHealthAndMinimumCustody(ctx, poolId); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventOpenConsolidate,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(existingMtp.Id), 10)),
		sdk.NewAttribute("owner", existingMtp.Address),
		sdk.NewAttribute("position", existingMtp.Position.String()),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(existingMtp.AmmPoolId), 10)),
		sdk.NewAttribute("collateral_asset", existingMtp.CollateralAsset),
		sdk.NewAttribute("collateral", existingMtp.Collateral.String()),
		sdk.NewAttribute("liabilities", existingMtp.Liabilities.String()),
		sdk.NewAttribute("custody", existingMtp.Custody.String()),
		sdk.NewAttribute("mtp_health", existingMtp.MtpHealth.String()),
		sdk.NewAttribute("stop_loss_price", existingMtp.StopLossPrice.String()),
		sdk.NewAttribute("take_profit_price", existingMtp.TakeProfitPrice.String()),
		sdk.NewAttribute("take_profit_borrow_factor", existingMtp.TakeProfitBorrowFactor.String()),
		sdk.NewAttribute("funding_fee_paid_custody", existingMtp.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_custody", existingMtp.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute("open_price", existingMtp.OpenPrice.String()),
	))

	return &types.MsgOpenResponse{
		Id: existingMtp.Id,
	}, nil
}
