package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

// Update params through gov proposal
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	// store params
	if err := k.SetParams(ctx, msg.Params); err != nil {
		return nil, err
	}

	if k.hooks != nil {
		pools := k.GetAllPools(ctx)
		for _, pool := range pools {
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
			if err != nil {
				return nil, fmt.Errorf("amm pool %d not found", pool.AmmPoolId)
			}

			err = k.hooks.AfterParamsChange(ctx, ammPool, pool)
			if err != nil {
				return nil, err
			}
		}
	}
	return &types.MsgUpdateParamsResponse{}, nil
}
