package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllLegacyPool(ctx)
	for _, pool := range pools {
		newPoolAssets := []types.PoolAsset{}
		for _, poolAsset := range pool.PoolAssets {
			newPoolAssets = append(newPoolAssets, types.PoolAsset{Token: poolAsset.Token, Weight: poolAsset.Weight, ExternalLiquidityRatio: sdk.OneDec()})
		}
		newPool := types.Pool{
			PoolId:            pool.PoolId,
			Address:           pool.Address,
			PoolParams:        pool.PoolParams,
			TotalShares:       pool.TotalShares,
			PoolAssets:        newPoolAssets,
			TotalWeight:       pool.TotalWeight,
			RebalanceTreasury: pool.RebalanceTreasury,
		}

		m.keeper.RemovePool(ctx, pool.PoolId)
		m.keeper.SetPool(ctx, newPool)
	}
	return nil
}
