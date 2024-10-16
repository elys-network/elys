package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) TakeOutCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool, amount math.Int) error {
	err := pool.UpdateCustody(ctx, mtp.CustodyAsset, amount, false, mtp.Position)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return nil
}
