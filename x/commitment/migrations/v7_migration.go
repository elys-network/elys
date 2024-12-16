package migrations

import (
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
		EnableClaim:             false,
	}
	m.keeper.SetParams(ctx, newParams)

	// Add missing wallet addresses to atom stakers DS
	for _, staker := range AtomStakers {
		m.keeper.SetAtomStaker(ctx, staker)
	}

	// Add kol addresses to kol DS
	for _, kol := range KolClaim {
		m.keeper.SetKol(ctx, kol)
	}

	return nil
}
