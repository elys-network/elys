package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) SetAssetInfo(goCtx context.Context, msg *types.MsgSetAssetInfo) (*types.MsgSetAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if msg.Creator != params.ModuleAdmin {
		return nil, types.ErrNotModuleAdmin
	}
	assetInfo := types.AssetInfo{
		Denom:      msg.Denom,
		Display:    msg.Display,
		BandTicker: msg.BandTicker,
		ElysTicker: msg.ElysTicker,
	}

	k.Keeper.SetAssetInfo(ctx, assetInfo)
	return &types.MsgSetAssetInfoResponse{}, nil
}

func (k msgServer) DeleteAssetInfo(goCtx context.Context, msg *types.MsgDeleteAssetInfo) (*types.MsgDeleteAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	if msg.Creator != params.ModuleAdmin {
		return nil, types.ErrNotModuleAdmin
	}
	k.RemoveAssetInfo(ctx, msg.Denom)
	return &types.MsgDeleteAssetInfoResponse{}, nil
}
