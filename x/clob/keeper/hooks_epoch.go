package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/elys-network/elys/v6/x/epochs/types"
)

// EpochHooks wrapper struct for incentive keeper
type EpochHooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = EpochHooks{}

// EpochHooks Return the wrapper struct
func (k Keeper) EpochHooks() EpochHooks {
	return EpochHooks{k}
}

// BeforeEpochStart implements EpochHooks
func (h EpochHooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	params := h.k.GetParams(ctx)
	if epochIdentifier == params.EpochIdentifier {
		allMarkets := h.k.GetAllPerpetualMarket(ctx)
		for _, market := range allMarkets {
			err := h.k.UpdateFundingRate(ctx, market)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// AfterEpochEnd implements EpochHooks
func (h EpochHooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) error {
	return nil
}
