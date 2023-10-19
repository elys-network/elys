package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) TakeOutCustody(ctx sdk.Context, mtp types.MTP, pool *types.Pool, custodyAsset string) error {
	_, custodyIndex := k.GetMTPAssetIndex(&mtp, "", custodyAsset)
	err := pool.UpdateBalance(ctx, mtp.CustodyAssets[custodyIndex], mtp.CustodyAmounts[custodyIndex], true)
	if err != nil {
		return err
	}

	err = pool.UpdateCustody(ctx, mtp.CustodyAssets[custodyIndex], mtp.CustodyAmounts[custodyIndex], false)
	if err != nil {
		return err
	}

	k.SetPool(ctx, *pool)

	return nil
}
