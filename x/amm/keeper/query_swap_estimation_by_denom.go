package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SwapEstimationByDenom(goCtx context.Context, req *types.QuerySwapEstimationByDenomRequest) (*types.QuerySwapEstimationByDenomResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	_ = baseCurrency

	inRoute, outRoute, amount, spotPrice, err := k.CalcSwapEstimationByDenom(ctx, req.Amount, req.DenomIn, req.DenomOut, baseCurrency)
	if err != nil {
		return nil, err
	}

	return &types.QuerySwapEstimationByDenomResponse{
		InRoute:   inRoute,
		OutRoute:  outRoute,
		Amount:    amount,
		SpotPrice: spotPrice,
	}, nil
}
