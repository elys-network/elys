package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MTP(goCtx context.Context, req *types.MTPRequest) (*types.MTPResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	creator, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return &types.MTPResponse{}, err
	}
	mtp, err := k.GetMTP(ctx, creator, req.Id)
	if err != nil {
		return &types.MTPResponse{}, err
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return &types.MTPResponse{}, status.Error(codes.NotFound, "base currency not found")
	}
	baseCurrency := entry.Denom

	mtpAndPrice, err := k.fillMTPData(ctx, mtp, baseCurrency)
	if err != nil {
		return &types.MTPResponse{}, err
	}

	return &types.MTPResponse{Mtp: mtpAndPrice}, nil
}
