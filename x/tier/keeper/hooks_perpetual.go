package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/v5/x/perpetual/types"
)

// Hooks wrapper struct for tvl keeper
type PerpetualHooks struct {
	k Keeper
}

var _ perpetualtypes.PerpetualHooks = PerpetualHooks{}

// Return the wrapper struct
func (k Keeper) PerpetualHooks() PerpetualHooks {
	return PerpetualHooks{k}
}

func (h PerpetualHooks) AfterParamsChange(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, EnableTakeProfitCustodyLiabilities bool) error {
	return nil
}

func (h PerpetualHooks) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (h PerpetualHooks) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (h PerpetualHooks) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}
