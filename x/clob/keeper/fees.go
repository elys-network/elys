package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) CollectTradingFees(ctx sdk.Context, market types.PerpetualMarket, trade types.Trade) error {
	buyerFeeRate := market.MakerFeeRate
	sellerFeeRate := market.TakerFeeRate
	if trade.IsBuyerTaker {
		buyerFeeRate = market.TakerFeeRate
		sellerFeeRate = market.MakerFeeRate
	}
	if trade.IsBuyerLiquidation {
		buyerFeeRate = math.LegacyZeroDec()
	}

	if trade.IsSellerLiquidation {
		sellerFeeRate = math.LegacyZeroDec()
	}
	tradeValue := trade.GetTradeValue()

	if buyerFeeRate.IsPositive() {
		buyerFees := tradeValue.Mul(buyerFeeRate).RoundInt()
		err := k.SendFromSubAccount(ctx, trade.BuyerSubAccount, market.GetInsuranceAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, buyerFees)))
		if err != nil {
			return err
		}
	}

	if sellerFeeRate.IsPositive() {
		sellerFees := tradeValue.Mul(sellerFeeRate).RoundInt()
		err := k.SendFromSubAccount(ctx, trade.SellerSubAccount, market.GetInsuranceAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, sellerFees)))
		if err != nil {
			return err
		}
	}

	return nil
}
