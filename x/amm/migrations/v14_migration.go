package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V14Migration(ctx sdk.Context) error {
	allPools := m.keeper.GetAllPool(ctx)
	for _, pool := range allPools {
		for i, poolAsset := range pool.PoolAssets {
			bankBalance := m.keeper.GetBankKeeper().GetBalance(ctx, sdk.MustAccAddressFromBech32(pool.Address), poolAsset.Token.Denom)
			pool.PoolAssets[i].Token = bankBalance
		}
		m.keeper.SetPool(ctx, pool)
	}
	return nil
}
