package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetWhitelist(goCtx context.Context, req *types.WhitelistRequest) (*types.WhitelistResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if req.Pagination != nil && req.Pagination.Limit > types.MaxPageLimit {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("page size greater than max %d", types.MaxPageLimit))
	}

	whitelist, page, err := k.GetWhitelistedAddress(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	whitelistAddressStrings := make([]string, len(whitelist))
	for i, whitelistAddress := range whitelist {
		whitelistAddressStrings[i] = whitelistAddress.String()
	}

	return &types.WhitelistResponse{
		Whitelist:  whitelistAddressStrings,
		Pagination: page,
	}, nil
}
