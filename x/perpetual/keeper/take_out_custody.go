package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) TakeOutCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool, amount sdk.Int) error {
	err := pool.UpdateBalance(ctx, mtp.CustodyAsset, amount, true, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, amount, false, mtp.Position)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return nil
}
