package keeper

//func (k Keeper) UpdateOraclePoolId(ctx sdk.Context, poolId uint64) error {
//	ammPool, found := k.GetPool(ctx, poolId)
//	if !found {
//		return types.ErrPoolNotFound
//	}
//	if !ammPool.PoolParams.UseOracle {
//		return nil
//	}
//	if len(ammPool.PoolAssets) > 2 {
//		return errors.New("cannot use more than two oracle pool assets")
//	}
//	baseDenom := ""
//	for _, poolAsset := range ammPool.PoolAssets {
//		assetInfo, found := k.oracleKeeper.GetAssetInfo(ctx, poolAsset.Token.Denom)
//		if !found {
//			return fmt.Errorf("asset info not found for pool asset %s", poolAsset.Token.Denom)
//		}
//		if assetInfo.Display != ptypes.USDC_DISPLAY {
//			baseDenom = assetInfo.Display
//		}
//	}
//	if baseDenom == "" {
//		return fmt.Errorf("invalid asset for oracle pool update id %d", poolId)
//	}
//	currencyPairProviders := k.oracleKeeper.CurrencyPairProviders(ctx)
//
//	for i, _ := range currencyPairProviders {
//		if currencyPairProviders[i].BaseDenom == baseDenom && currencyPairProviders[i].ExternLiquidityProvider != "" {
//			currencyPairProviders[i].PoolId = poolId
//		}
//	}
//	k.oracleKeeper.SetCurrencyPairProviders(ctx, currencyPairProviders)
//
//	return nil
//}
