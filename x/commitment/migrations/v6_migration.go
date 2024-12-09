package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	legacy := m.keeper.GetLegacyParams(ctx)
	newParams := types.Params{
		VestingInfos:           legacy.VestingInfos,
		TotalCommitted:         legacy.TotalCommitted,
		NumberOfCommitments:    legacy.NumberOfCommitments,
		EnableVestNow:          legacy.EnableVestNow,
		StartAtomStakersHeight: 0,
		StartCadetsHeight:      0,
		StartGovernorsHeight:   0,
		StartNftHoldersHeight:  0,
		EndAtomStakersHeight:   0,
		EndCadetsHeight:        0,
		EndGovernorsHeight:     0,
		EndNftHoldersHeight:    0,
	}
	m.keeper.SetParams(ctx, newParams)

	return nil
}
