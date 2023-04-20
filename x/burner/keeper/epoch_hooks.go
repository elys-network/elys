package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/elys-network/elys/x/epochs/types"
)

// BeforeEpochStart performs a no-op
func (k Keeper) BeforeEpochStart(_ sdk.Context, _ string, _ int64) {}

// AfterEpochEnd distributes vested tokens at the end of each epoch
func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, _ int64) {

	// Future Improvement: check all VestingInfos and get all VestingTokens by denom
	// 	so we can iterate different denoms in different EpochIdentifiers
	params := k.GetParams(ctx)
	if epochIdentifier == params.EpochIdentifier {
		denoms := k.bankKeeper.GetAllDenomMetaData(ctx)
		for _, denom := range denoms {
			k.Logger(ctx).Info("Burning tokens for denom", denom.Base)
			if err := k.BurnTokens(ctx, denom.Base); err != nil {
				k.Logger(ctx).Error("Error burning tokens", "denom", denom.Base)
				panic(err)
			}
		}
	}
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for commitments keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// BeforeEpochStart implements EpochHooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

// AfterEpochEnd implements EpochHooks
func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
