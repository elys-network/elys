package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/accountedpool/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	allLegacyAccountedPool := m.keeper.GetAllLegacyAccountedPool(ctx)

	for _, legacyAccountedPool := range allLegacyAccountedPool {
		accountedPool := types.AccountedPool{
			PoolId:           legacyAccountedPool.PoolId,
			TotalTokens:      []sdk.Coin{},
			NonAmmPoolTokens: legacyAccountedPool.NonAmmPoolTokens,
		}

		for _, asset := range legacyAccountedPool.PoolAssets {
			accountedPool.TotalTokens = append(accountedPool.TotalTokens, asset.Token)
		}

		m.keeper.SetAccountedPool(ctx, accountedPool)

	}
	return nil
}
