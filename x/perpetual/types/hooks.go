package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v5/x/amm/types"
)

type PerpetualHooks interface {
	AfterParamsChange(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, EnableTakeProfitCustodyLiabilities bool) error
	// AfterPerpetualPositionOpen is called after OpenLong or OpenShort position.
	// This should be used to update pool health
	AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error

	// AfterPerpetualPositionModified is called after a position gets modified.
	// This should be used to update pool health
	AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error

	// AfterPerpetualPositionClosed is called after a position gets closed.
	// This should be used to update pool health
	AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error
}

var _ PerpetualHooks = MultiPerpetualHooks{}

// combine multiple perpetual hooks, all hook functions are run in array sequence.
type MultiPerpetualHooks []PerpetualHooks

// Creates hooks for the Amm Module.
func NewMultiPerpetualHooks(hooks ...PerpetualHooks) MultiPerpetualHooks {
	return hooks
}

func (h MultiPerpetualHooks) AfterParamsChange(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, EnableTakeProfitCustodyLiabilities bool) error {
	for i := range h {
		err := h[i].AfterParamsChange(ctx, ammPool, perpetualPool, EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiPerpetualHooks) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	for i := range h {
		err := h[i].AfterPerpetualPositionOpen(ctx, ammPool, perpetualPool, sender, EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiPerpetualHooks) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	for i := range h {
		err := h[i].AfterPerpetualPositionModified(ctx, ammPool, perpetualPool, sender, EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h MultiPerpetualHooks) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool Pool, sender sdk.AccAddress, EnableTakeProfitCustodyLiabilities bool) error {
	for i := range h {
		err := h[i].AfterPerpetualPositionClosed(ctx, ammPool, perpetualPool, sender, EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return err
		}
	}
	return nil
}
