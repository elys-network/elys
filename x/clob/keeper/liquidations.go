package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) LiquidationClose(ctx sdk.Context, perpetual types.Perpetual, market types.PerpetualMarket, closer sdk.AccAddress) error {
	msg := types.MsgPlaceMarketOrder{
		Creator:      perpetual.Owner,
		MarketId:     perpetual.MarketId,
		BaseQuantity: perpetual.Quantity.Abs(),
		OrderType:    types.OrderType_ORDER_TYPE_MARKET_SELL,
	}
	if perpetual.IsShort() {
		msg.OrderType = types.OrderType_ORDER_TYPE_MARKET_BUY
		orderFilled, err := k.ExecuteMarketBuyOrder(ctx, market, msg)
		if err != nil {
			return err
		}
		if !orderFilled {
			return fmt.Errorf("unable to liquidate %d for market %d, order cannot be filled", perpetual.Id, perpetual.MarketId)
		}
	} else {
		orderFilled, err := k.ExecuteMarketSellOrder(ctx, market, msg)
		if err != nil {
			return err
		}
		if !orderFilled {
			return fmt.Errorf("unable to liquidate %d for market %d, order cannot be filled", perpetual.Id, perpetual.MarketId)
		}
	}
	// perpetual has been deleted from KV store and balances have been reverted
	closerAmountDec, err := market.LiquidationFeeShareRate.Mul(utils.IntToDec(perpetual.Margin))
	if err != nil {
		return err
	}
	closerAmount, err := closerAmountDec.SdkIntTrim()
	if err != nil {
		return err
	}
	subAccount, err := k.GetSubAccount(ctx, perpetual.GetOwnerAccAddress(), perpetual.MarketId)
	if err != nil {
		return err
	}
	err = k.SendFromSubAccount(ctx, subAccount, closer, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, closerAmount)))
	if err != nil {
		return err
	}
	return nil
}
