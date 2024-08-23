package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CommittedTokensLocked(goCtx context.Context, req *types.QueryCommittedTokensLockedRequest) (*types.QueryCommittedTokensLockedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)

	if err != nil {
		return nil, err
	}

	listPositionAndInterest, _, err := k.GetPositionsForAddress(ctx, address, nil)

	if err != nil {
		return nil, err
	}

	totalLocked, totalCommitted := sdk.Coins{}, sdk.Coins{}
	for _, positionAndInterest := range listPositionAndInterest {

		commitments := k.commKeeper.GetCommitments(ctx, positionAndInterest.Position.Position.GetPositionAddress())
		tl, tc := commitments.CommittedTokensLocked(ctx)

		totalLocked = totalLocked.Add(tl...)
		totalCommitted = totalCommitted.Add(tc...)

	}

	return &types.QueryCommittedTokensLockedResponse{
		Address:         address.String(),
		TotalCommitted:  totalCommitted,
		LockedCommitted: totalLocked,
	}, nil
}
