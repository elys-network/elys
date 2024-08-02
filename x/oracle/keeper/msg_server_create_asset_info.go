package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (k msgServer) CreateAssetInfo(goCtx context.Context, msg *types.MsgCreateAssetInfo) (*types.MsgCreateAssetInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgCreateAssetInfoResponse{}, nil
}
