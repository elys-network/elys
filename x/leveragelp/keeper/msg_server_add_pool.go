package keeper

import (
	"context"
	"fmt"

	sdkmath "cosmossdk.io/math"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/leveragelp/types"
)

func (k msgServer) AddPool(goCtx context.Context, msg *types.MsgAddPool) (*types.MsgAddPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	ammPool, ammFound := k.amm.GetPool(ctx, msg.Pool.AmmPoolId)

	if !ammFound {
		return nil, fmt.Errorf("amm pool not found")
	}

	if !ammPool.PoolParams.UseOracle {
		return nil, fmt.Errorf("amm pool does not use oracle")
	}

	_, found := k.GetPool(ctx, msg.Pool.AmmPoolId)
	if found {
		return nil, fmt.Errorf("pool already exists")
	}

	maxLeverageAllowed := k.GetMaxLeverageParam(ctx)
	leverage := sdkmath.LegacyMinDec(msg.Pool.LeverageMax, maxLeverageAllowed)

	newPool := types.NewPool(ammPool.PoolId, leverage)
	for _, asset := range ammPool.PoolAssets {
		newPool.AssetLeverageAmounts = append(newPool.AssetLeverageAmounts, &types.AssetLeverageAmount{
			Denom:           asset.Token.Denom,
			LeveragedAmount: sdkmath.ZeroInt(),
		})
	}
	k.SetPool(ctx, newPool)

	if k.hooks != nil {
		err := k.hooks.AfterEnablingPool(ctx, ammPool)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgAddPoolResponse{}, nil
}
