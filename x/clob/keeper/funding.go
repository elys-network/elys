package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
)

func (k Keeper) SettleFunding(ctx sdk.Context, subAccount *types.SubAccount, market types.PerpetualMarket, perpetual *types.Perpetual) error {
	currentFundingRate := k.GetFundingRate(ctx, market.Id)
	fundingRateApplied := currentFundingRate.Rate.Sub(perpetual.EntryFundingRate)
	if !fundingRateApplied.IsZero() {
		paymentSign := int64(-1)
		if perpetual.IsShort() {
			paymentSign = 1
		}
		twapPrice := k.GetCurrentTwapPrice(ctx, market.Id)
		if !twapPrice.IsZero() {
			fundingPnL := fundingRateApplied.Mul(perpetual.Quantity.Abs().Mul(twapPrice)).RoundInt().MulRaw(paymentSign)
			if fundingPnL.IsPositive() {
				err := k.AddToSubAccount(ctx, market.GetAccount(), *subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, fundingPnL)))
				if err != nil {
					return err
				}
			} else {
				err := k.SendFromSubAccount(ctx, *subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, fundingPnL.Abs())))
				if err != nil {
					return err
				}
			}
		}
	}

	perpetual.EntryFundingRate = currentFundingRate.Rate
	return nil
}
