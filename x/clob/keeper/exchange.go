package keeper

import (
	"cosmossdk.io/math"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) Exchange(ctx sdk.Context, trade types.Trade) error {
	if trade.Quantity.LTE(math.LegacyZeroDec()) {
		return errors.New("trade quantity must be greater than zero")
	}
	if trade.SellerSubAccount.IsIsolated() && trade.SellerSubAccount.Id != trade.MarketId {
		return errors.New("trade market id and subAccounts market id does not match for seller")
	}
	if trade.BuyerSubAccount.IsIsolated() && trade.BuyerSubAccount.Id != trade.MarketId {
		return errors.New("trade market id and subAccounts market id does not match for buyer")
	}

	market, err := k.GetPerpetualMarket(ctx, trade.MarketId)
	if err != nil {
		return err
	}

	currentFundingRate := k.GetFundingRate(ctx, market.Id)

	buyerPositionBefore := math.LegacyZeroDec()
	sellerPositionBefore := math.LegacyZeroDec()

	// Buyer Changes
	var buyerPerpetual types.Perpetual

	buyerPerpetualOwner, buyerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.BuyerSubAccount.GetOwnerAccAddress(), trade.BuyerSubAccount.Id)
	if buyerAlreadyOwn {
		buyerPerpetual, err = k.GetPerpetual(ctx, trade.MarketId, buyerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}

		buyerPositionBefore = buyerPerpetual.Quantity

		err = k.SettleFunding(ctx, &trade.BuyerSubAccount, market, &buyerPerpetual)
		if err != nil {
			return err
		}

	}

	buyerPerpetual, err = k.SettleMarginAndRPnL(ctx, market, buyerPerpetual, trade.IsBuyerLiquidation, trade, true)
	if err != nil {
		return err
	}

	if buyerPerpetual.IsZero() {
		k.DeletePerpetual(ctx, buyerPerpetual)
		k.DeletePerpetualOwner(ctx, buyerPerpetual.GetOwnerAccAddress(), trade.BuyerSubAccount.Id)
	} else {
		buyerPerpetual.EntryFundingRate = currentFundingRate.Rate
		k.SetPerpetual(ctx, buyerPerpetual)

		if !buyerAlreadyOwn {
			buyerPerpetualOwner = types.PerpetualOwner{
				Owner:        buyerPerpetual.Owner,
				SubAccountId: trade.BuyerSubAccount.Id,
				MarketId:     trade.MarketId,
				PerpetualId:  buyerPerpetual.Id,
			}
			k.SetPerpetualOwner(ctx, buyerPerpetualOwner)
		}
	}

	// Seller Changes
	var sellerPerpetual types.Perpetual

	sellerPerpetualOwner, sellerAlreadyOwn := k.GetPerpetualOwner(ctx, trade.SellerSubAccount.GetOwnerAccAddress(), trade.SellerSubAccount.Id)
	if sellerAlreadyOwn {
		sellerPerpetual, err = k.GetPerpetual(ctx, trade.MarketId, sellerPerpetualOwner.PerpetualId)
		if err != nil {
			return err
		}

		sellerPositionBefore = sellerPerpetual.Quantity

		err = k.SettleFunding(ctx, &trade.SellerSubAccount, market, &sellerPerpetual)
		if err != nil {
			return err
		}

	}

	sellerPerpetual, err = k.SettleMarginAndRPnL(ctx, market, sellerPerpetual, trade.IsSellerLiquidation, trade, false)
	if err != nil {
		return err
	}
	if sellerPerpetual.IsZero() {
		k.DeletePerpetual(ctx, sellerPerpetual)
		k.DeletePerpetualOwner(ctx, sellerPerpetual.GetOwnerAccAddress(), trade.SellerSubAccount.Id)
	} else {
		sellerPerpetual.EntryFundingRate = currentFundingRate.Rate
		k.SetPerpetual(ctx, sellerPerpetual)

		if !sellerAlreadyOwn {
			sellerPerpetualOwner = types.PerpetualOwner{
				Owner:        sellerPerpetual.Owner,
				SubAccountId: trade.SellerSubAccount.Id,
				MarketId:     trade.MarketId,
				PerpetualId:  sellerPerpetual.Id,
			}
			k.SetPerpetualOwner(ctx, sellerPerpetualOwner)
		}
	}

	err = k.CollectTradingFees(ctx, market, trade)
	if err != nil {
		return err
	}
	// Market Changes
	market.UpdateTotalOpenInterest(buyerPositionBefore, sellerPositionBefore, trade.Quantity)
	k.SetPerpetualMarket(ctx, market)
	err = k.SetTwapPrices(ctx, trade)
	if err != nil {
		return err
	}
	return nil
}
