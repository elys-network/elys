package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) Params(goCtx context.Context, request *types.ParamsRequest) (*types.ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.ParamsResponse{Params: k.GetParams(ctx)}, nil
}
