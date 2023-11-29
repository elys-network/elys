package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) TakeOutCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool, custodyAsset string) error {
	err := pool.UpdateBalance(ctx, custodyAsset, mtp.Custodies.AmountOf(custodyAsset), true, mtp.Position)
	if err != nil {
		return err
	}

	err = pool.UpdateCustody(ctx, custodyAsset, mtp.Custodies.AmountOf(custodyAsset), false, mtp.Position)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return nil
}
