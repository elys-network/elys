package keeper

import (
	"context"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
	"strconv"
)

func (k msgServer) AddCollateral(goCtx context.Context, msg *types.MsgAddCollateral) (*types.MsgAddCollateralResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}

	mtp, err := k.Keeper.GetMTP(ctx, msg.PoolId, creator, msg.Id)
	if err != nil {
		return nil, err
	}

	initialCollateralCoin := sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral)
	initialCustody := mtp.Custody
	initialLiabilities := mtp.Liabilities

	pool, found := k.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return nil, errors.New("pool not found")
	}

	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return nil, err
	}

	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, closingRatio, err := k.MTPTriggerChecksAndUpdates(ctx, &mtp, &pool, &ammPool)
	if err != nil {
		return nil, err
	}

	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	if forceClosed {
		usdcPrice, err := k.GetUSDCPrice(ctx)
		if err != nil {
			return nil, err
		}
		k.EmitForceClose(ctx, "add_collateral", mtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, msg.Creator, allInterestsPaid, tradingAssetPrice, totalPerpetualFeesCoins, closingPrice, initialCollateralCoin, initialCustody, initialLiabilities, usdcPrice, closingRatio)
		// hooks are being called inside MTPTriggerChecksAndUpdates
		return &types.MsgAddCollateralResponse{}, nil
	}

	// trigger check happened just above
	finalCollateralCoin, err := k.Keeper.AddCollateral(ctx, &mtp, &pool, msg.AddCollateral, &ammPool, true)
	if err != nil {
		return nil, err
	}

	mtp.PartialLiquidationDone = false
	err = k.Keeper.SetMTP(ctx, &mtp)
	if err != nil {
		return nil, err
	}

	perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd := k.GetPerpFeesInUSD(ctx, totalPerpetualFeesCoins)
	interestAmtInUSD := k.amm.CalculateUSDValue(ctx, mtp.CustodyAsset, interestAmt).Dec()

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAddCollateral,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("owner", mtp.Address),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("collateral_added", finalCollateralCoin.String()),
		sdk.NewAttribute("updated_collateral", mtp.Collateral.String()),
		sdk.NewAttribute("updated_liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("funding_fee_amount", fundingFeeAmt.String()),
		sdk.NewAttribute("funding_amount_distributed", fundingAmtDistributed.String()),
		sdk.NewAttribute("interest_amount", interestAmt.String()),
		sdk.NewAttribute("interest_amount_in_usd", interestAmtInUSD.String()),
		sdk.NewAttribute("insurance_amount", insuranceAmt.String()),
		sdk.NewAttribute("funding_fee_paid_custody", mtp.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_custody", mtp.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute("trading_asset_price", tradingAssetPrice.String()),
		sdk.NewAttribute(types.AttributeKeyPerpFee, perpFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeySlippage, slippageFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeyWeightBreakingFee, weightBreakingFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeTakerFees, takerFeesInUsd.String()),
	))

	return &types.MsgAddCollateralResponse{}, nil
}
