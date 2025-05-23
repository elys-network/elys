package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/leveragelp/types"
)

// Update enabled pools through gov proposal
func (k msgServer) UpdateEnabledPools(goCtx context.Context, msg *types.MsgUpdateEnabledPools) (*types.MsgUpdateEnabledPoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// store params
	params := k.GetParams(ctx)
	params.EnabledPools = msg.EnabledPools
	if err := k.SetParams(ctx, &params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateEnabledPoolsResponse{}, nil
}
