package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) SettleFunding(ctx sdk.Context, subAccount *types.SubAccount, market types.PerpetualMarket, perpetual *types.Perpetual) error {
	currentFundingRate := k.GetFundingRate(ctx, market.Id)
	fundingRateApplied, err := currentFundingRate.Rate.Sub(perpetual.EntryFundingRate)
	if err != nil {
		return err
	}
	if !fundingRateApplied.IsZero() {
		paymentSign := int64(-1)
		if perpetual.IsShort() {
			paymentSign = 1
		}
		fundingPnLDec, err := utils.IntToDec(perpetual.Quantity.Abs().MulRaw(paymentSign)).Mul(fundingRateApplied)
		if err != nil {
			return err
		}
		fundingPnL, err := fundingPnLDec.SdkIntTrim()
		if err != nil {
			return err
		}
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

	perpetual.EntryFundingRate = currentFundingRate.Rate
	return nil
}
