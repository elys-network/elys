package keeper

import (
	"context"
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

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	if mtp.Position == types.Position_LONG && mtp.CollateralAsset == baseCurrency {
		if msg.AddCollateral.Denom != mtp.CollateralAsset {
			return nil, errors.New("denom not same as collateral asset")
		}

		if msg.AddCollateral.Denom != mtp.LiabilitiesAsset {
			return nil, errors.New("denom not same as liabilities asset")
		}

		pool, found := k.GetPool(ctx, mtp.AmmPoolId)
		if !found {
			return nil, errors.New("pool not found")
		}

		ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
		if err != nil {
			return nil, err
		}

		repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, err := k.MTPTriggerChecksAndUpdates(ctx, &mtp, &pool, &ammPool)
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
			k.EmitForceClose(ctx, "add_collateral", mtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, msg.Creator, allInterestsPaid, tradingAssetPrice, totalPerpetualFeesCoins, closingPrice, initialCollateralCoin, initialCustody, initialLiabilities, usdcPrice)
			// hooks are being called inside MTPTriggerChecksAndUpdates
			return &types.MsgAddCollateralResponse{}, nil
		}

		// interest amount has been paid from custody
		params := k.GetParams(ctx)
		maxAmount := mtp.Liabilities.Sub(params.LongMinimumLiabilityAmount)
		if !maxAmount.IsPositive() {
			return nil, fmt.Errorf("cannot reduce liabilties less than %s", params.LongMinimumLiabilityAmount.String())
		}

		var finalAmount math.Int
		if msg.AddCollateral.Amount.LT(maxAmount) {
			finalAmount = msg.AddCollateral.Amount
		} else {
			finalAmount = maxAmount
		}

		mtp.Liabilities = mtp.Liabilities.Sub(finalAmount)
		err = pool.UpdateLiabilities(mtp.LiabilitiesAsset, finalAmount, false, mtp.Position)
		if err != nil {
			return nil, err
		}

		mtp.Collateral = mtp.Collateral.Add(finalAmount)
		err = pool.UpdateCollateral(mtp.CollateralAsset, finalAmount, true, mtp.Position)
		if err != nil {
			return nil, err
		}

		finalCollateralCoin := sdk.NewCoin(msg.AddCollateral.Denom, finalAmount)
		err = k.SendToAmmPool(ctx, creator, &ammPool, sdk.NewCoins(finalCollateralCoin))
		if err != nil {
			return nil, err
		}

		err = k.SetMTP(ctx, &mtp)
		if err != nil {
			return nil, err
		}
		k.SetPool(ctx, pool)

		if k.hooks != nil {
			err = k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator)
			if err != nil {
				return nil, err
			}
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

	} else {
		msgOpen := types.MsgOpen{
			Creator:         msg.Creator,
			Position:        mtp.Position,
			Leverage:        math.LegacyZeroDec(),
			Collateral:      msg.AddCollateral,
			TakeProfitPrice: math.LegacyZeroDec(),
			StopLossPrice:   math.LegacyZeroDec(),
			PoolId:          mtp.AmmPoolId,
		}
		if err = msgOpen.ValidateBasic(); err != nil {
			return nil, err
		}
		_, err = k.Open(goCtx, &msgOpen)
		if err != nil {
			return nil, err
		}

	}
	return &types.MsgAddCollateralResponse{}, nil
}
