package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	params := m.keeper.GetLegacyParams(ctx)
	m.keeper.SetParams(ctx, params)

	priceFeeders := m.keeper.GetAllLegacyPriceFeeder(ctx)
	for _, priceFeeder := range priceFeeders {
		m.keeper.SetPriceFeeder(ctx, priceFeeder)
		m.keeper.RemoveLegacyPriceFeeder(ctx, priceFeeder.Feeder)
	}
	return nil
}
