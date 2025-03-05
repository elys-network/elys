package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetParams(ctx, data.Params)
	for _, v := range data.SubAccounts {
		k.SetSubAccount(ctx, v)
	}
	for _, v := range data.PerpetualMarkets {
		k.SetPerpetualMarket(ctx, v)
	}
	for _, v := range data.Perpetuals {
		k.SetPerpetual(ctx, v)
	}
	for _, v := range data.OrderBooks {
		k.SetPerpetualOrder(ctx, *v)
	}
	for _, v := range data.PerpetualOwners {
		k.SetPerpetualOwner(ctx, v)
	}
	for _, v := range data.MarketPrices {
		k.SetMarketPrice(ctx, v.MarketId, v.LastPrice, false)
		k.SetMarketPrice(ctx, v.MarketId, v.LastPrice, true)
	}
	for _, v := range data.PerpetualCounters {
		k.setPerpetualCounter(ctx, *v)
	}
	for _, v := range data.Escrows {
		k.SetEscrow(ctx, *v)
	}
}
