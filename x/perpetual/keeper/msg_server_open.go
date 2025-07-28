package keeper

import (
	"context"
	"fmt"
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	enabledPools := k.GetParams(ctx).EnabledPools
	found := slices.Contains(enabledPools, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolNotEnabled, fmt.Sprintf("poolId: %d", msg.PoolId))
	}

	return k.Keeper.Open(ctx, msg)
}
