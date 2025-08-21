package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V22Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.MinimumFundingRate = math.LegacyMustNewDecFromStr("0.1")
	params.SecondLiquidationTriggerRatio = math.LegacyMustNewDecFromStr("0.67")
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}
	return nil
}
