package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) TakeOutCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool, custodyAsset string) error {
	_, custodyIndex := k.GetMTPAssetIndex(&mtp, "", custodyAsset)
	err := pool.UpdateBalance(ctx, mtp.Custodies[custodyIndex].Denom, mtp.Custodies[custodyIndex].Amount, true)
	if err != nil {
		return err
	}

	err = pool.UpdateCustody(ctx, mtp.Custodies[custodyIndex].Denom, mtp.Custodies[custodyIndex].Amount, false)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return nil
}
