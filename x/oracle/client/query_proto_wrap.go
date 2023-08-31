package client

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/keeper"
	"github.com/elys-network/elys/x/oracle/types"
)

type Querier struct {
	K keeper.Keeper
}

func NewQuerier(k keeper.Keeper) Querier {
	return Querier{K: k}
}

func (q Querier) PriceAll(ctx sdk.Context, req types.QueryAllPriceRequest) (*types.QueryAllPriceResponse, error) {
	return q.K.PriceAll(ctx, &req)
}
