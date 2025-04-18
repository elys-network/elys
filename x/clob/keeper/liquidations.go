package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) LiquidationClose(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket, bot sdk.AccAddress) error {
	botSubAccount, err := k.GetSubAccount(ctx, bot, market.Id)
	if err != nil {
		return err
	}
	subAccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), market.Id)
	if err != nil {
		return err
	}
	liquidationPrice, err := k.GetLiquidationPrice(ctx, perpetual, market)
	if err != nil {
		return err
	}
	trade := types.Trade{
		MarketId: perpetual.MarketId,
		Quantity: perpetual.Quantity.Abs(),
		Price:    liquidationPrice,
	}
	if perpetual.IsShort() {
		trade.BuyerSubAccount = subAccount
		trade.SellerSubAccount = botSubAccount
	} else {
		trade.BuyerSubAccount = botSubAccount
		trade.SellerSubAccount = subAccount
	}

	// TODO Shoyld bot liquidations impact TWAP price? I think no
	err = k.Exchange(ctx, trade)
	if err != nil {
		return err
	}

	// perpetual has been deleted from KV store and balances have been reverted
	botAmount := market.LiquidationFeeShareRate.MulInt(perpetual.Margin).TruncateInt()
	err = k.SendFromSubAccount(ctx, subAccount, bot, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, botAmount)))
	if err != nil {
		return err
	}
	return nil
}
