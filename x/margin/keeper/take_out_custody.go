package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) TakeOutCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool) error {
	err := pool.UpdateBalance(ctx, mtp.CustodyAsset, mtp.CustodyAmount, true)
	if err != nil {
		return err
	}

	err = pool.UpdateCustody(ctx, mtp.CustodyAsset, mtp.CustodyAmount, false)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return nil
}
