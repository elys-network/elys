package keeper

import (
	"context"

	"github.com/elys-network/elys/x/masterchef/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) AddExternalIncentive(goCtx context.Context, msg *types.MsgAddExternalIncentive) (*types.MsgAddExternalIncentiveResponse, error) {
	return &types.MsgAddExternalIncentiveResponse{}, nil
}

func (k msgServer) ClaimRewards(goCtx context.Context, msg *types.MsgClaimRewards) (*types.MsgClaimRewardsResponse, error) {
	return &types.MsgClaimRewardsResponse{}, nil
}
