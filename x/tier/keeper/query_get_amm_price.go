package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetAmmPrice(goCtx context.Context, req *types.QueryGetAmmPriceRequest) (*types.QueryGetAmmPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	resp := k.amm.CalcAmmPrice(ctx, req.Denom, uint64(req.Decimal))

	return &types.QueryGetAmmPriceResponse{
		Result: &types.GetAmmPriceResponseResult{
			Total: resp,
		},
	}, nil
}
