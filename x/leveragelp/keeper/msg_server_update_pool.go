package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

// Update max leverage for a pool through gov proposal
func (k msgServer) UpdatePool(goCtx context.Context, msg *types.MsgUpdatePool) (*types.MsgUpdatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, fmt.Errorf("pool does not exists for pool id %d", msg.PoolId)
	}

	maxLeverageAllowed := k.GetMaxLeverageParam(ctx)
	if maxLeverageAllowed.LT(msg.LeverageMax) {
		return nil, fmt.Errorf("max leverage allowed is less than the leverage max")
	}

	pool.LeverageMax = msg.LeverageMax
	pool.MaxLeveragelpRatio = msg.MaxLeveragelpRatio

	k.SetPool(ctx, pool)

	return &types.MsgUpdatePoolResponse{}, nil
}
