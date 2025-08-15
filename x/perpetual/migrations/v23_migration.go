package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V23Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.SecondLiquidationTriggerRatio = math.LegacyMustNewDecFromStr("0.67")
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	for _, pool := range m.keeper.GetAllPools(ctx) {
		pool.MtpSafetyFactor = params.SafetyFactor
		m.keeper.SetPool(ctx, pool)
	}
	return nil
}
