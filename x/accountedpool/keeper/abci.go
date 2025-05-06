package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {

	// TODO this is being used for price feeder, migrate to query in price feeder and the remove this
	//accountedPools := k.GetAllAccountedPool(ctx)
	//for _, accountedPool := range accountedPools {
	//	oracleAccountedPool := oracletypes.AccountedPool{
	//		PoolId:      accountedPool.PoolId,
	//		TotalTokens: accountedPool.TotalTokens,
	//	}
	//
	//	k.oracleKeeper.SetAccountedPool(ctx, oracleAccountedPool)
	//}
}
