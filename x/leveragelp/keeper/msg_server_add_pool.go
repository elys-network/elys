package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (k msgServer) AddPool(goCtx context.Context, msg *types.MsgAddPool) (*types.MsgAddPoolResponse, error) {
	var newPool types.Pool
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	pool, ammFound := k.amm.GetPool(ctx, msg.Pool.AmmPoolId)

	if ammFound && pool.PoolParams.UseOracle {
		_, found := k.GetPool(ctx, msg.Pool.AmmPoolId)

		if !found {
			maxLeverageAllowed := k.GetMaxLeverageParam(ctx)
			leverage := sdk.MinDec(msg.Pool.LeverageMax, maxLeverageAllowed)

			newPool.AmmPoolId = msg.Pool.AmmPoolId
			newPool.Closed = msg.Pool.Closed
			newPool.Enabled = msg.Pool.Enabled
			newPool.LeverageMax = leverage
			newPool.Health = sdk.NewDec(0)
			newPool.LeveragedLpAmount = sdk.NewInt(0)
			k.SetPool(ctx, newPool)
		}
	}

	return &types.MsgAddPoolResponse{}, nil
}
