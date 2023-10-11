package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, mtp *types.MTP, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {

	poolId := mtp.AmmPoolId
	pool, found := k.OpenLongChecker.GetPool(ctx, poolId)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	if !k.OpenLongChecker.IsPoolEnabled(ctx, poolId) {
		return nil, sdkerrors.Wrap(types.ErrMTPDisabled, fmt.Sprintf("poolId: %d", poolId))
	}

	ammPool, err := k.OpenLongChecker.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	mtp, err = k.OpenConsolidateLong(ctx, poolId, mtp, msg)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(k.GenerateOpenEvent(mtp))

	if k.hooks != nil {
		k.hooks.AfterLeveragelpPositionModified(ctx, ammPool, pool)
	}

	return &types.MsgOpenResponse{}, nil
}
