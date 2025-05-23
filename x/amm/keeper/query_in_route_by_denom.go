package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InRouteByDenom(goCtx context.Context, req *types.QueryInRouteByDenomRequest) (*types.QueryInRouteByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return nil, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	inRoute, err := k.CalcInRouteByDenom(ctx, req.DenomIn, req.DenomOut, baseCurrency)
	if err != nil {
		return nil, err
	}

	return &types.QueryInRouteByDenomResponse{
		InRoute: inRoute,
	}, nil
}
