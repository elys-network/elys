package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OutRouteByDenom(goCtx context.Context, req *types.QueryOutRouteByDenomRequest) (*types.QueryOutRouteByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	outRoute, err := k.CalcOutRouteByDenom(ctx, req.DenomOut, req.DenomIn, baseCurrency)
	if err != nil {
		return nil, err
	}

	return &types.QueryOutRouteByDenomResponse{
		OutRoute: outRoute,
	}, nil
}
