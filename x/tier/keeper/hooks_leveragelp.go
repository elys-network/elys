package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	leveragelptypes "github.com/elys-network/elys/v4/x/leveragelp/types"
)

type LeverageLpHooks struct {
	k Keeper
}

var _ leveragelptypes.LeverageLpHooks = LeverageLpHooks{}

// Return the wrapper struct
func (k Keeper) LeverageLpHooks() LeverageLpHooks {
	return LeverageLpHooks{k}
}

func (h LeverageLpHooks) AfterEnablingPool(_ sdk.Context, _ ammtypes.Pool) error {
	return nil
}

func (h LeverageLpHooks) AfterDisablingPool(_ sdk.Context, _ ammtypes.Pool) error {
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress, _ ammtypes.Pool) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress, _ ammtypes.Pool) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress, _ ammtypes.Pool) error {
	h.k.RetrieveAllPortfolio(ctx, sender)
	return nil
}
