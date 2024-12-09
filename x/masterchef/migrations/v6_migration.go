package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/masterchef/types"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	stableStakePoolId := uint64(stabletypes.PoolId)

	// Set pool info for stablestake pool as it was not set earlier
	m.keeper.SetPoolInfo(ctx, types.PoolInfo{
		PoolId: stableStakePoolId,
		// reward wallet address
		RewardWallet: ammtypes.NewPoolRevenueAddress(stableStakePoolId).String(),
		// multiplier for lp rewards
		Multiplier: math.LegacyNewDec(1),
		// Eden APR, updated at every distribution
		EdenApr: math.LegacyZeroDec(),
		// Dex APR, updated at every distribution
		DexApr: math.LegacyZeroDec(),
		// Gas APR, updated at every distribution
		GasApr: math.LegacyZeroDec(),
		// External Incentive APR, updated at every distribution
		ExternalIncentiveApr: math.LegacyZeroDec(),
		// external reward denoms on the pool
		ExternalRewardDenoms: []string{},
		EnableEdenRewards:    true,
	})
	return nil
}
