package keeper

import (
	"context"

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
	_, shares, _, weightBalanceBonus, err := k.amm.JoinPoolEst(ctx, req.AmmPoolId, sdk.Coins{leverageCoin})
	if err != nil {
		return nil, err
	}
	params := k.stableKeeper.GetParams(ctx)

	return &types.QueryOpenEstResponse{
		PositionSize:       shares,
		WeightBalanceRatio: weightBalanceBonus,
		BorrowFee:          params.InterestRate,
	}, nil
}

func (k Keeper) CloseEst(goCtx context.Context, req *types.QueryCloseEstRequest) (*types.QueryCloseEstResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := sdk.MustAccAddressFromBech32(req.Owner)
	position, err := k.GetPosition(ctx, owner, req.Id)
	if err != nil {
		return nil, err
	}

	exitCoins, weightBalanceBonus, err := k.amm.ExitPoolEst(ctx, position.AmmPoolId, req.LpAmount, position.Collateral.Denom)
	if err != nil {
		return nil, err
	}

	// Repay with interest
	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress(), position.BorrowPoolId)

	// Ensure position.LeveragedLpAmount is not zero to avoid division by zero
	if position.LeveragedLpAmount.IsZero() {
		return nil, types.ErrAmountTooLow
	}

	repayAmount := debt.GetTotalLiablities().Mul(req.LpAmount).Quo(position.LeveragedLpAmount)
	userAmount := exitCoins[0].Amount.Sub(repayAmount)

	return &types.QueryCloseEstResponse{
		Liability:          repayAmount,
		WeightBalanceRatio: weightBalanceBonus,
		AmountReturned:     userAmount,
	}, nil
}
