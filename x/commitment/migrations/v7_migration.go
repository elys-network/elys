package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/commitment/types"
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
		StartKolClaimHeight:     legacy.StartAtomStakersHeight,
		EndKolClaimHeight:       legacy.EndAtomStakersHeight,
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

	// For testnet only
	if ctx.ChainID() == "elysicstestnet-1" {
		addresses := []string{"elys1u8c28343vvhwgwhf29w6hlcz73hvq7lwxmrl46", "elys1va67r6h8y489kz89rgpqtlysrv6pd6hckr3pe4", "elys1zu4tahdxl0d585wpnu03nne4wkleyz5qdllh60"}
		for _, address := range addresses {
			m.keeper.SetCadet(ctx, types.Cadet{
				Address: address,
				Amount:  math.NewInt(500000000),
			})
			m.keeper.SetGovernor(ctx, types.Governor{
				Address: address,
				Amount:  math.NewInt(500000000),
			})
			m.keeper.SetNFTHodler(ctx, types.NftHolder{
				Address: address,
				Amount:  math.NewInt(500000000),
			})

			m.keeper.SetKol(ctx, types.KolList{
				Address:  address,
				Amount:   math.NewInt(17143000000),
				Claimed:  false,
				Refunded: false,
			})

			m.keeper.SetAtomStaker(ctx, types.AtomStaker{
				Address: address,
				Amount:  math.NewInt(30000000),
			})
		}
	}

	return nil
}
