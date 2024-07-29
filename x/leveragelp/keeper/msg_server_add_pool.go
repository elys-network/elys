package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) AddPools(goCtx context.Context, msg *types.MsgAddPool) (*types.MsgAddPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	pool, found := k.amm.GetPool(ctx, msg.Pool.AmmPoolId)

	if found && pool.PoolParams.UseOracle {
		_, found := k.GetPool(ctx, msg.Pool.AmmPoolId)

		if !found {
			k.SetPool(ctx, msg.Pool)
		}
	}

	return &types.MsgAddPoolResponse{}, nil
}
