package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) CreateAssetInfo(goCtx context.Context, msg *types.MsgCreateAssetInfo) (*types.MsgCreateAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetAssetInfo(ctx, msg.Denom)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "asset already set")
	}

	assetInfo := types.AssetInfo{
		Denom:         msg.Denom,
		Display:       msg.Display,
		BandTicker:    msg.BandTicker,
		BinanceTicker: msg.BinanceTicker,
		OsmosisTicker: msg.OsmosisTicker,
	}

	k.SetAssetInfo(ctx, assetInfo)
	return &types.MsgCreateAssetInfoResponse{}, nil
}

func (k msgServer) UpdateAssetInfo(goCtx context.Context, msg *types.MsgUpdateAssetInfo) (*types.MsgUpdateAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	assetInfo := types.AssetInfo{
		Denom:         msg.Denom,
		Display:       msg.Display,
		BandTicker:    msg.BandTicker,
		BinanceTicker: msg.BinanceTicker,
		OsmosisTicker: msg.OsmosisTicker,
	}

	k.SetAssetInfo(ctx, assetInfo)
	return &types.MsgUpdateAssetInfoResponse{}, nil
}

func (k msgServer) DeleteAssetInfo(goCtx context.Context, msg *types.MsgDeleteAssetInfo) (*types.MsgDeleteAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.RemoveAssetInfo(ctx, msg.Denom)
	return &types.MsgDeleteAssetInfoResponse{}, nil
}
