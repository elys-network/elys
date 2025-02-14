package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OpenEst(goCtx context.Context, req *types.QueryOpenEstRequest) (*types.QueryOpenEstResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	leveragedAmount := req.Leverage.MulInt(req.CollateralAmount).TruncateInt()
	leverageCoin := sdk.NewCoin(req.CollateralAsset, leveragedAmount)
	_, shares, slippage, weightBalanceBonus, swapFee, err := k.amm.JoinPoolEst(ctx, req.AmmPoolId, sdk.Coins{leverageCoin})
	if err != nil {
		return nil, err
	}
	params := k.stableKeeper.GetParams(ctx)

	return &types.QueryOpenEstResponse{
		PositionSize:       shares,
		WeightBalanceRatio: weightBalanceBonus,
		BorrowFee:          params.InterestRate,
		Slippage:           slippage,
		SwapFee:            swapFee,
	}, nil
}

func (k Keeper) CloseEst(goCtx context.Context, req *types.QueryCloseEstRequest) (*types.QueryCloseEstResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx, _ := sdk.UnwrapSDKContext(goCtx).CacheContext()
	owner := sdk.MustAccAddressFromBech32(req.Owner)
	position, err := k.GetPosition(ctx, owner, req.Id)
	if err != nil {
		return nil, err
	}
	if req.LpAmount.GT(position.LeveragedLpAmount) {
		return nil, errors.New("request lp amount is greater than position lp amount")
	}
	pool, found := k.GetPool(ctx, position.AmmPoolId)
	if !found {
		return nil, errors.New("leverage lp pool not found")
	}

	closingRatio := req.LpAmount.ToLegacyDec().Quo(position.LeveragedLpAmount.ToLegacyDec())
	finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, _, weightBreakingFee, exitSlippageFee, swapFee, err := k.CheckHealthStopLossThenRepayAndClose(ctx, &position, &pool, closingRatio, false)
	if err != nil {
		return nil, err
	}

	return &types.QueryCloseEstResponse{
		RepayAmount:       repayAmount,
		FinalClosingRatio: finalClosingRatio,
		ClosingLpAmount:   totalLpAmountToClose,
		CoinsToAmm:        coinsForAmm,
		UserReturnTokens:  userReturnTokens,
		ExitWeightFee:     exitFeeOnClosingPosition,
		WeightBreakingFee: weightBreakingFee,
		ExitSlippageFee:   exitSlippageFee,
		ExitSwapFee:       swapFee,
	}, nil
}
