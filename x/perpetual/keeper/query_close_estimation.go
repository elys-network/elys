package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) CloseEstimation(goCtx context.Context, req *types.QueryCloseEstimationRequest) (*types.QueryCloseEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}
	mtp, err := k.CloseLongChecker.GetMTP(ctx, address, req.PositionId)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	// Retrieve Pool
	pool, found := k.CloseLongChecker.GetPool(ctx, mtp.AmmPoolId)
	if !found {
		return &types.QueryCloseEstimationResponse{}, errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
	}

	_ = pool

	// Retrieve AmmPool
	ammPool, err := k.CloseLongChecker.GetAmmPool(ctx, mtp.AmmPoolId, mtp.CustodyAsset)
	if err != nil {
		return &types.QueryCloseEstimationResponse{}, err
	}

	_ = ammPool

	// TODO: estimate settle borrow interest

	// TODO: estimate take out custody

	// TODO: estimate swap and repay

	return &types.QueryCloseEstimationResponse{
		Position:     mtp.Position,
		PositionSize: sdk.NewCoin(mtp.CustodyAsset, mtp.Custody),
		Liabilities:  sdk.NewCoin(mtp.LiabilitiesAsset, mtp.Liabilities),
		PriceImpact:  sdk.ZeroDec(),
		SwapFee:      sdk.ZeroDec(),
		RepayAmount:  sdk.Coin{},
	}, nil
}
