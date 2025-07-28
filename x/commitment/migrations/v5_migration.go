package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/commitment/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	params := m.keeper.GetLegacyParams(ctx)
	newParams := types.Params{
		VestingInfos:        params.VestingInfos,
		TotalCommitted:      params.TotalCommitted,
		NumberOfCommitments: params.NumberOfCommitments,
		EnableVestNow:       false,
	}
	m.keeper.SetParams(ctx, newParams)

	return nil
}
