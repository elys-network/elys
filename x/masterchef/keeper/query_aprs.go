package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Aprs(goCtx context.Context, req *types.QueryAprsRequest) (*types.QueryAprsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	usdcAprUsdc, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_USDC_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprUsdc, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_USDC_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	usdcAprEdenb, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDENB_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprEdenb, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDENB_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	usdcAprEden, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDEN_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprEden, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDEN_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	edenbAprEden, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_EDEN_PROGRAM, Denom: ptypes.EdenB})
	if err != nil {
		return nil, err
	}

	usdcAprElys, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_ELYS_PROGRAM, Denom: ptypes.BaseCurrency})
	if err != nil {
		return nil, err
	}

	edenAprElys, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_ELYS_PROGRAM, Denom: ptypes.Eden})
	if err != nil {
		return nil, err
	}

	edenbAprElys, err := k.CalculateApr(ctx, &types.QueryAprRequest{WithdrawType: commitmenttypes.EarnType_ELYS_PROGRAM, Denom: ptypes.EdenB})
	if err != nil {
		return nil, err
	}

	return &types.QueryAprsResponse{
		UsdcAprUsdc:  usdcAprUsdc.Dec(),
		EdenAprUsdc:  edenAprUsdc.Dec(),
		UsdcAprEdenb: usdcAprEdenb.Dec(),
		EdenAprEdenb: edenAprEdenb.Dec(),
		UsdcAprEden:  usdcAprEden.Dec(),
		EdenAprEden:  edenAprEden.Dec(),
		EdenbAprEden: edenbAprEden.Dec(),
		UsdcAprElys:  usdcAprElys.Dec(),
		EdenAprElys:  edenAprElys.Dec(),
		EdenbAprElys: edenbAprElys.Dec(),
	}, nil
}
