package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/leveragelp/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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
	_, shares, slippage, weightBalanceBonus, swapFee, takerFees, weightRewardAmount, err := k.amm.JoinPoolEst(ctx, req.AmmPoolId, sdk.Coins{leverageCoin})
	if err != nil {
		return nil, err
	}
	pool, found := k.stableKeeper.GetPoolByDenom(ctx, req.CollateralAsset)
	if !found {
		return nil, errors.New("borrow pool not found")
	}

	return &types.QueryOpenEstResponse{
		PositionSize:              shares,
		WeightBalanceRatio:        weightBalanceBonus.Dec(),
		BorrowFee:                 pool.InterestRate,
		Slippage:                  slippage.Dec(),
		SwapFee:                   swapFee.Dec(),
		TakerFee:                  takerFees.Dec(),
		WeightBalanceRewardAmount: weightRewardAmount,
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

	closingRatio := osmomath.BigDecFromSDKInt(req.LpAmount).Quo(position.GetBigDecLeveragedLpAmount())
	finalClosingRatio, totalLpAmountToClose, coinsForAmm, repayAmount, userReturnTokens, exitFeeOnClosingPosition, _, weightBreakingFee, exitSlippageFee, swapFee, takerFee, err := k.CheckHealthStopLossThenRepayAndClose(ctx, &position, &pool, closingRatio, false)
	if err != nil {
		return nil, err
	}

	return &types.QueryCloseEstResponse{
		RepayAmount:       repayAmount,
		FinalClosingRatio: finalClosingRatio.Dec(),
		ClosingLpAmount:   totalLpAmountToClose,
		CoinsToAmm:        coinsForAmm,
		UserReturnTokens:  userReturnTokens,
		ExitWeightFee:     exitFeeOnClosingPosition.Dec(),
		WeightBreakingFee: weightBreakingFee.Dec(),
		ExitSlippageFee:   exitSlippageFee.Dec(),
		ExitSwapFee:       swapFee.Dec(),
		ExitTakerFee:      takerFee.Dec(),
	}, nil
}
