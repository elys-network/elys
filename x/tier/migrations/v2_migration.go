package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	year, month, day := ctx.BlockTime().Date()
	dateToday := time.Date(year, month, day, 0, 0, 0, 0, ctx.BlockTime().Location())
	for i := -8; i <= 0; i++ {
		date := dateToday.AddDate(0, 0, i)
		portfolios := m.keeper.GetLegacyPortfolios(ctx, date.String())
		for _, portfolio := range portfolios {
			m.keeper.SetPortfolio(ctx, date.String(), portfolio)
			m.keeper.RemoveLegacyPortfolio(ctx, date.String(), portfolio.Creator)
		}
	}
	return nil
}
