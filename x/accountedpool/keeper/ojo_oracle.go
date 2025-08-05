package keeper

// UpdateOjoOracleAccountedPool TODO this is being used for price feeder, migrate to query in price feeder and the remove this
// We cannot use AccountedPool inside ojo/oracle due to import cycle issue
//func (k Keeper) UpdateOjoOracleAccountedPool(ctx sdk.Context) {
//	accountedPools := k.GetAllAccountedPool(ctx)
//	for _, accountedPool := range accountedPools {
//		oracleAccountedPool := oracletypes.AccountedPool{
//			PoolId:      accountedPool.PoolId,
//			TotalTokens: accountedPool.TotalTokens,
//		}
//
//		k.oracleKeeper.SetAccountedPool(ctx, oracleAccountedPool)
//	}
//}
