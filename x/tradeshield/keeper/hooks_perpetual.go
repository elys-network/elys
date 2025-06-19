package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/v6/x/perpetual/types"
)

type PerpetualHooks struct {
	k Keeper
}

var _ perpetualtypes.PerpetualHooks = PerpetualHooks{}

// Return the wrapper struct
func (k Keeper) PerpetualHooks() PerpetualHooks {
	return PerpetualHooks{k}
}

func (h PerpetualHooks) AfterParamsChange(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) error {
	return nil
}

func (h PerpetualHooks) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress) error {
	return nil
}

func (h PerpetualHooks) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress) error {
	return nil
}

func (h PerpetualHooks) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress) error {
	return nil
}

func (h PerpetualHooks) AfterPositionDestroyed(ctx sdk.Context, owner sdk.AccAddress, poolId uint64, positionId uint64) error {
	_, err := h.k.DeletePendingPerpetualOrdersForAddressAndPool(ctx, owner, poolId, positionId)
	if err != nil {
		return err
	}
	return nil
}
