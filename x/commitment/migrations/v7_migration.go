package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	legacy := m.keeper.GetLegacyParams(ctx)
	newParams := types.Params{
		VestingInfos:            legacy.VestingInfos,
		TotalCommitted:          legacy.TotalCommitted,
		NumberOfCommitments:     legacy.NumberOfCommitments,
		EnableVestNow:           legacy.EnableVestNow,
		StartAirdropClaimHeight: legacy.StartAtomStakersHeight,
		EndAirdropClaimHeight:   legacy.EndAtomStakersHeight,
		TotalElysClaimed:        math.ZeroInt(),
		TotalEdenClaimed:        math.ZeroInt(),
	}
	m.keeper.SetParams(ctx, newParams)

	return nil
}
