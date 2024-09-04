package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	for i := -8; i <= 0; i++ {
		date := m.keeper.GetDateAfterDaysFromContext(ctx, i)
		portfolios := m.keeper.GetLegacyPortfolios(ctx, date)
		for _, portfolio := range portfolios {
			m.keeper.SetPortfolio(ctx, portfolio)
			m.keeper.RemoveLegacyPortfolio(ctx, date, portfolio.Creator)
		}
	}
	return nil
}
