package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k msgServer) EnablePool(goCtx context.Context, msg *types.MsgEnablePool) (*types.MsgEnablePoolResponse, error) {
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetPool(ctx, msg.PoolId)
	if found {
		return nil, fmt.Errorf("pool already exists: %d", msg.PoolId)
	}

	leverageLpPool, found := k.LeverageLpKeeper.GetPool(ctx, msg.PoolId)
	if !found || !leverageLpPool.Enabled || leverageLpPool.Closed {
		return nil, fmt.Errorf("leverage lp not enabled for the pool %d", msg.PoolId)
	}

	ammPool, err := k.GetAmmPool(ctx, msg.PoolId)
	if err != nil {
		return nil, fmt.Errorf("amm pool not found %d", msg.PoolId)
	}

	pool = types.NewPool(msg.PoolId)
	err = pool.InitiatePool(ctx, &ammPool)
	if err != nil {
		return nil, err
	}

	k.SetPool(ctx, pool)

	if k.hooks != nil {
		err = k.hooks.AfterEnablingPool(ctx, ammPool)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgEnablePoolResponse{}, nil
}

func (k msgServer) DisablePool(goCtx context.Context, msg *types.MsgDisablePool) (*types.MsgDisablePoolResponse, error) {
	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, fmt.Errorf("pool not found: %d", msg.PoolId)
	}

	ammPool, err := k.GetAmmPool(ctx, msg.PoolId)
	if err != nil {
		return nil, fmt.Errorf("amm pool not found %d", msg.PoolId)
	}

	// TODO Figure out a way to close positions before deleting it

	k.RemovePool(ctx, msg.PoolId)

	if k.hooks != nil {
		err = k.hooks.AfterDisablingPool(ctx, ammPool)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgDisablePoolResponse{}, nil
}
