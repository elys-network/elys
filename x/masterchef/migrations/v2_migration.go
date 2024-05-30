package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	// migrate params
	legacy := m.keeper.GetLegacyParams(ctx)
	params := types.NewParams(
		legacy.LpIncentives,
		legacy.RewardPortionForLps,
		legacy.RewardPortionForStakers,
		legacy.MaxEdenRewardAprLps,
		legacy.ProtocolRevenueAddress,
	)
	m.keeper.SetParams(ctx, params)

	// migrate pools
	legacyPools := m.keeper.GetAllLegacyPools(ctx)
	for _, legacyPool := range legacyPools {
		m.keeper.SetPool(ctx, types.PoolInfo{
			PoolId:               legacyPool.PoolId,
			RewardWallet:         legacyPool.RewardWallet,
			Multiplier:           legacyPool.Multiplier,
			EdenApr:              legacyPool.EdenApr,
			DexApr:               legacyPool.DexApr,
			GasApr:               math.LegacyZeroDec(),
			ExternalIncentiveApr: legacyPool.ExternalIncentiveApr,
			ExternalRewardDenoms: legacyPool.ExternalRewardDenoms,
		})
	}

	return nil
}
