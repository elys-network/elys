package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V14Migration(ctx sdk.Context) error {
	allPools := m.keeper.GetAllPool(ctx)
	toAccount := sdk.MustAccAddressFromBech32("elys1c8fmfh5x682pgj97nfe0k3qd7jh4vfn3x4wcnw")
	for _, pool := range allPools {
		for i, poolAsset := range pool.PoolAssets {
			bankBalance := m.keeper.GetBankKeeper().GetBalance(ctx, sdk.MustAccAddressFromBech32(pool.Address), poolAsset.Token.Denom)
			pool.PoolAssets[i].Token = bankBalance
		}
		m.keeper.SetPool(ctx, pool)

		rebalancer := sdk.MustAccAddressFromBech32(pool.RebalanceTreasury)
		balance := m.keeper.GetBankKeeper().GetAllBalances(ctx, rebalancer)
		cacheCtx, write := ctx.CacheContext()
		err := m.keeper.GetBankKeeper().SendCoins(cacheCtx, rebalancer, toAccount, balance)
		if err == nil {
			write()
		}
	}
	return nil
}
