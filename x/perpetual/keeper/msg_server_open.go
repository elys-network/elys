package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) Open(goCtx context.Context, msg *types.MsgOpen) (*types.MsgOpenResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	enabledPools := k.GetParams(ctx).EnabledPools
	found := false
	for _, poolId := range enabledPools {
		if poolId == msg.PoolId {
			found = true
			break
		}
	}
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolNotEnabled, fmt.Sprintf("poolId: %d", msg.PoolId))
	}

	return k.Keeper.Open(ctx, msg)
}
