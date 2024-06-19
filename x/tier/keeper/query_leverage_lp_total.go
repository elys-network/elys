package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LeverageLpTotal(goCtx context.Context, req *types.QueryLeverageLpTotalRequest) (*types.QueryLeverageLpTotalResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(req.User)
	total := k.RetreiveLeverageLpTotal(ctx, sender)

	return &types.QueryLeverageLpTotalResponse{
		Total: total,
	}, nil
}
