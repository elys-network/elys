package keeper

import (
	"context"
	"cosmossdk.io/math"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
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

	positions, pageResponse, err := k.GetPositionsForAddress(ctx, address, req.Pagination)
	if err != nil {
		return nil, err
	}

	var positionCommitedTokens []types.PositionCommitedToken

	for _, position := range positions {

		commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress())
		tl, tc := commitments.CommittedTokensLocked(ctx)
		if len(tl) == 0 {
			tl = sdk.NewCoins(sdk.NewCoin(ammtypes.GetPoolShareDenom(position.AmmPoolId), math.ZeroInt()))
		}
		if len(tc) == 0 {
			tc = sdk.NewCoins(sdk.NewCoin(ammtypes.GetPoolShareDenom(position.AmmPoolId), math.ZeroInt()))
		}

		positionCommitedToken := types.PositionCommitedToken{
			AmmPoolId:       position.AmmPoolId,
			PositionId:      position.Id,
			BorrowPoolId:    position.BorrowPoolId,
			CollateralDenom: position.Collateral.Denom,
			LockedCommitted: tl[0],
			TotalCommitted:  tc[0],
		}
		positionCommitedTokens = append(positionCommitedTokens, positionCommitedToken)

	}

	return &types.QueryCommittedTokensLockedResponse{
		Address:               address.String(),
		PositionCommitedToken: positionCommitedTokens,
		Pagination:            pageResponse,
	}, nil
}
