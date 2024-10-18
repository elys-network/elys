package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
)

type LeverageLpHooks struct {
	k Keeper
}

var _ leveragelptypes.LeverageLpHooks = LeverageLpHooks{}

// Return the wrapper struct
func (k Keeper) LeverageLpHooks() LeverageLpHooks {
	return LeverageLpHooks{k}
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}
