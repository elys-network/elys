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

func (k Keeper) QueryPositions(goCtx context.Context, req *types.PositionsRequest) (*types.PositionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	if req.Pagination != nil && req.Pagination.Limit > types.MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	positions, page, err := k.GetPositions(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	updatedLeveragePositions := []*types.QueryPosition{}
	for i, position := range positions {
		pool, found := k.amm.GetPool(ctx, position.AmmPoolId)
		if !found {
			return nil, errorsmod.Wrap(ammtypes.ErrPoolNotFound, fmt.Sprintf("poolId: %d", position.AmmPoolId))
		}
		lp_price, err := pool.LpTokenPrice(ctx, k.oracleKeeper)
		if err != nil {
			return nil, err
		}
		
		lp_usd_price := position.LeveragedLpAmount.Mul(lp_price.TruncateInt())
		price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Collateral.Denom)
		updated_leverage :=  lp_usd_price.Quo(lp_usd_price.Sub(position.Liabilities.Mul(price.TruncateInt())))
		
		updatedLeveragePositions[i] = &types.QueryPosition{
			Position:        position,
			UpdatedLeverage: updated_leverage,
		}
	}

	return &types.PositionsResponse{
		Positions:  updatedLeveragePositions,
		Pagination: page,
	}, nil
}
