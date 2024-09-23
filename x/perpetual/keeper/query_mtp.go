package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
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

	info, found := k.oracleKeeper.GetAssetInfo(ctx, mtp.TradingAsset)
	if !found {
		return nil, fmt.Errorf("asset not found")
	}
	trading_asset_price, found := k.oracleKeeper.GetAssetPrice(ctx, info.Display)
	if !found {
		return nil, fmt.Errorf("asset price not found")
	}

	return &types.MTPResponse{Mtp: &types.MtpAndPrice{Mtp: &mtp, TradingAssetPrice: trading_asset_price.Price}}, nil
}
