package keeper

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) error {
	allMarkets := k.GetAllPerpetualMarket(ctx)

	for _, market := range allMarkets {
		marketPrice := k.GetLastMarketPrice(ctx, market.Id)

		{
			fmt.Println("Buy Book: ")
			buyIterator := k.GetBuyOrderIterator(ctx, market.Id)

			defer buyIterator.Close()
			var list []types.PerpetualOrder
			for ; buyIterator.Valid(); buyIterator.Next() {
				var val types.PerpetualOrder
				k.cdc.MustUnmarshal(buyIterator.Value(), &val)
				list = append(list, val)
			}
			fmt.Println(list)
		}

		{
			fmt.Println("Sell Book: ")
			buyIterator := k.GetSellOrderIterator(ctx, market.Id)

			defer buyIterator.Close()
			var list []types.PerpetualOrder
			for ; buyIterator.Valid(); buyIterator.Next() {
				var val types.PerpetualOrder
				k.cdc.MustUnmarshal(buyIterator.Value(), &val)
				list = append(list, val)
			}
			fmt.Println(list)
		}

		err := k.ExecuteMarket(ctx, market.Id)
		if err != nil {
			return err
		}

		marketPrice = k.GetLastMarketPrice(ctx, market.Id)
		if !marketPrice.IsZero() {
			k.SetLastMarketPrice(ctx, market.Id, marketPrice, true)
		}
	}
	return nil
}
