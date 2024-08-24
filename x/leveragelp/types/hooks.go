package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LeverageLpHooks interface {
	// AfterLeverageLpPositionOpen is called after Open position.
	AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress) error

	// AfterLeverageLpPositionClose is called after a position gets closed.
	AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress) error

	// AfterLeverageLpPositionConsolidate is called after a position gets closed.
	AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress) error
}

var _ LeverageLpHooks = MultiLeverageLpHooks{}

// combine multiple leverageLp hooks, all hook functions are run in array sequence.
type MultiLeverageLpHooks []LeverageLpHooks


func NewMultiLeverageLpHooks(hooks ...LeverageLpHooks) MultiLeverageLpHooks {
	return hooks
}

func (h MultiLeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress) error {
	for i := range h {
		err := h[i].AfterLeverageLpPositionOpen(ctx, sender)
		if err != nil {
			return err
		}
	}
	return nil
}


func (h MultiLeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress) error {
	for i := range h {
		err := h[i].AfterLeverageLpPositionClose(ctx, sender)
		if err != nil {
			return err
		}
	}
	return nil
}


func (h MultiLeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress) error {
	for i := range h {
		err := h[i].AfterLeverageLpPositionOpenConsolidate(ctx, sender)
		if err != nil {
			return err
		}
	}
	return nil
}
