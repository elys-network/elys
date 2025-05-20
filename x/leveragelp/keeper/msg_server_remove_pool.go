package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
)

func (k msgServer) RemovePool(goCtx context.Context, msg *types.MsgRemovePool) (*types.MsgRemovePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}
	ammPool, found := k.amm.GetPool(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolDoesNotExist, "amm pool (%d) not found", msg.Id)
	}

	pool, found := k.GetPool(ctx, msg.Id)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool with id %d not found", msg.Id)
	}

	if pool.LeveragedLpAmount.IsPositive() {
		return nil, errorsmod.Wrap(types.ErrPoolLeverageAmountNotZero, pool.LeveragedLpAmount.String())
	}

	k.DeletePool(ctx, msg.Id)

	if k.hooks != nil {
		err := k.hooks.AfterDisablingPool(ctx, ammPool)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgRemovePoolResponse{}, nil
}
