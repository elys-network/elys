package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	oldParams := m.keeper.GetLegacyParams(ctx)
	params := types.DefaultParams()
	params.TotalValue = oldParams.TotalValue
	params.InterestRate = oldParams.InterestRate
	m.keeper.SetParams(ctx, params)
	return nil
}
