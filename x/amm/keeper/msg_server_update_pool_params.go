package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k msgServer) UpdatePoolParams(goCtx context.Context, msg *types.MsgUpdatePoolParams) (*types.MsgUpdatePoolParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Sender {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Sender)
	}

	poolId, poolParams, err := k.Keeper.UpdatePoolParams(ctx, msg.PoolId, *msg.PoolParams)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdatePoolParamsResponse{
		PoolId:     poolId,
		PoolParams: &poolParams,
	}, nil
}
