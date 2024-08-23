package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) AfterLeverageLpPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, sender sdk.AccAddress) error {
	k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (k Keeper) AfterLeverageLpPositionClose(ctx sdk.Context, ammPool ammtypes.Pool, sender sdk.AccAddress) error {
	k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

func (k Keeper) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, ammPool ammtypes.Pool, sender sdk.AccAddress) error {
	k.RetrieveAllPortfolio(ctx, sender)
	return nil
}

type LeverageLpHooks struct {
	k Keeper
}

var _ leveragelptypes.LeverageLpHooks = LeverageLpHooks{}

// Return the wrapper struct
func (k Keeper) LeverageLpHooks() LeverageLpHooks {
	return LeverageLpHooks{k}
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, sender sdk.AccAddress) error {
	return h.k.AfterLeverageLpPositionOpen(ctx, ammPool, sender)
}

func (h LeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, ammPool ammtypes.Pool, sender sdk.AccAddress) error {
	return h.k.AfterLeverageLpPositionClose(ctx, ammPool, sender)
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, ammPool ammtypes.Pool, sender sdk.AccAddress) error {
	return h.k.AfterLeverageLpPositionOpenConsolidate(ctx, ammPool, sender)
}