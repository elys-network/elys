package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) CalcExitPoolCoinsFromShares(ctx sdk.Context, exitingShares math.Int, exitFee sdk.Dec) (exitedCoins sdk.Coins, err error) {
	return CalcExitPool(ctx, *p, exitingShares, exitFee)
}
