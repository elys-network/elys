package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V22Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.NumberPerBlock = 5
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	allPools := m.keeper.GetAllPools(ctx)
	for _, pool := range allPools {
		pool.AdlTriggerRatio = params.PoolMaxLiabilitiesThreshold.Add(math.LegacyMustNewDecFromStr("0.02"))
		m.keeper.SetPool(ctx, pool)
	}
	return nil
}
