package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
)

type LeverageLpHooks interface {
	AfterEnablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error

	AfterDisablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error

	// AfterLeverageLpPositionOpen is called after Open position.
	AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error

	// AfterLeverageLpPositionClose is called after a position gets closed.
	AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error

	// AfterLeverageLpPositionConsolidate is called after a position gets closed.
	AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error
}

var _ LeverageLpHooks = MultiLeverageLpHooks{}

// combine multiple leverageLp hooks, all hook functions are run in array sequence.
type MultiLeverageLpHooks []LeverageLpHooks

func NewMultiLeverageLpHooks(hooks ...LeverageLpHooks) MultiLeverageLpHooks {
	return hooks
}

func (h MultiLeverageLpHooks) AfterEnablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	for i := range h {
		err := h[i].AfterEnablingPool(ctx, ammPool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiLeverageLpHooks) AfterDisablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	for i := range h {
		err := h[i].AfterDisablingPool(ctx, ammPool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiLeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	for i := range h {
		err := h[i].AfterLeverageLpPositionOpen(ctx, sender, ammPool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiLeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	for i := range h {
		err := h[i].AfterLeverageLpPositionClose(ctx, sender, ammPool)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiLeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	for i := range h {
		err := h[i].AfterLeverageLpPositionOpenConsolidate(ctx, sender, ammPool)
		if err != nil {
			return err
		}
	}
	return nil
}
