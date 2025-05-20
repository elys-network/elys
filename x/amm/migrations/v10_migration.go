package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/amm/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)
	params := types.Params{
		PoolCreationFee:                  legacyParams.PoolCreationFee,
		SlippageTrackDuration:            legacyParams.SlippageTrackDuration,
		BaseAssets:                       legacyParams.BaseAssets,
		WeightBreakingFeeExponent:        legacyParams.WeightBreakingFeeExponent,
		WeightBreakingFeeMultiplier:      legacyParams.WeightBreakingFeeMultiplier,
		WeightBreakingFeePortion:         legacyParams.WeightBreakingFeePortion,
		WeightRecoveryFeePortion:         legacyParams.WeightRecoveryFeePortion,
		ThresholdWeightDifference:        legacyParams.ThresholdWeightDifference,
		AllowedPoolCreators:              legacyParams.AllowedPoolCreators,
		ThresholdWeightDifferenceSwapFee: math.LegacyMustNewDecFromStr("0.125"),
	}

	m.keeper.SetParams(ctx, params)
	return nil
}
