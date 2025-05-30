package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (m Migrator) V17Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	newParams := types.Params{
		LeverageMax:         legacyParams.LeverageMax,
		MaxOpenPositions:    legacyParams.MaxOpenPositions,
		PoolOpenThreshold:   legacyParams.PoolOpenThreshold,
		SafetyFactor:        legacyParams.SafetyFactor,
		WhitelistingEnabled: legacyParams.WhitelistingEnabled,
		EpochLength:         legacyParams.EpochLength,
		FallbackEnabled:     legacyParams.FallbackEnabled,
		NumberPerBlock:      legacyParams.NumberPerBlock,
		EnabledPools:        []uint64{},
	}

	err := m.keeper.SetParams(ctx, &newParams)
	if err != nil {
		return err
	}
	return nil
}
