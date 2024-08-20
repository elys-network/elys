package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPositionsByPool(goCtx context.Context, req *types.PositionsByPoolRequest) (*types.PositionsByPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	positions, pageRes, err := k.GetPositionsForPool(ctx, req.AmmPoolId, req.Pagination)
	if err != nil {
		return nil, err
	}

	pool, found := k.amm.GetPool(ctx, req.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(ammtypes.ErrPoolNotFound, fmt.Sprintf("poolId: %d", req.AmmPoolId))
	}

	lp_price, err := pool.LpTokenPrice(ctx, k.oracleKeeper)
	if err != nil {
		return nil, err
	}

	updatedLeveragePositions := []*types.QueryPosition{}
	for i ,position := range positions {
		lp_usd_price := position.LeveragedLpAmount.Mul(lp_price.TruncateInt())
		price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Collateral.Denom)
		updated_leverage :=  lp_usd_price.Quo(lp_usd_price.Sub(position.Liabilities.Mul(price.TruncateInt())))
		updatedLeveragePositions[i] = &types.QueryPosition{
			Position: position,
			UpdatedLeverage: updated_leverage,
		}
	}

	return &types.PositionsByPoolResponse{
		Positions:  updatedLeveragePositions,
		Pagination: pageRes,
	}, nil
}
