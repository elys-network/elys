package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/assetprofile/types"
)

func (k msgServer) AddEntry(goCtx context.Context, msg *types.MsgAddEntry) (*types.MsgAddEntryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgAddEntryResponse{}, nil
}
