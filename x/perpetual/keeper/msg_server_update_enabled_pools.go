package keeper

import (
	"context"
	"slices"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

// Update enabled pools through gov proposal
func (k msgServer) UpdateEnabledPools(goCtx context.Context, msg *types.MsgUpdateEnabledPools) (*types.MsgUpdateEnabledPoolsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// store params
	params := k.GetParams(ctx)
	for _, add := range msg.AddPools {
		if !slices.Contains(params.EnabledPools, add) {
			params.EnabledPools = append(params.EnabledPools, add)
		}
	}

	params.EnabledPools = slices.DeleteFunc(params.EnabledPools, func(element uint64) bool {
		return slices.Contains(msg.RemovePools, element)
	})

	slices.Sort(params.EnabledPools)
	if err := k.SetParams(ctx, &params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateEnabledPoolsResponse{}, nil
}
