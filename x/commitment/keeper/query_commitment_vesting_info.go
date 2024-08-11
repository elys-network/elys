package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CommitmentVestingInfo(goCtx context.Context, req *types.QueryCommitmentVestingInfoRequest) (*types.QueryCommitmentVestingInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	address := sdk.MustAccAddressFromBech32(req.Address)
	commitment := k.GetCommitments(ctx, address)
	vestingTokens := commitment.GetVestingTokens()

	totalVesting := sdk.ZeroInt()
	vestingDetails := make([]types.VestingDetails, 0)
	for i, vesting := range vestingTokens {
		vestingDetail := types.VestingDetails{
			Id:              fmt.Sprintf("%d", i),
			TotalVesting:    vesting.TotalAmount,
			Claimed:         vesting.ClaimedAmount,
			VestedSoFar:     vesting.VestedSoFar(ctx),
			RemainingBlocks: vesting.NumBlocks - (ctx.BlockHeight() - vesting.StartBlock),
		}

		vestingDetails = append(vestingDetails, vestingDetail)
		totalVesting = totalVesting.Add(vesting.TotalAmount.Sub(vesting.ClaimedAmount))
	}

	return &types.QueryCommitmentVestingInfoResponse{
		Total:          totalVesting,
		VestingDetails: vestingDetails,
	}, nil
}
