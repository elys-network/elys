package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) isElysRoutedMultihop(_ sdk.Context, route types.MultihopRoute, inDenom, outDenom string) (isRouted bool) {
	if route.Length() != 2 {
		return false
	}
	intemediateDenoms := route.IntermediateDenoms()
	if len(intemediateDenoms) != 1 {
		return false
	}
	if inDenom == outDenom {
		return false
	}

	poolIds := route.PoolIds()

	return poolIds[0] != poolIds[1]
}
