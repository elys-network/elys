package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Position(goCtx context.Context, req *types.PositionRequest) (*types.PositionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	position, err := k.GetPosition(ctx, req.Address, req.Id)
	if err != nil {
		return nil, err
	}

	commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress().String())
	totalLocked, _ := commitments.CommittedTokensLocked(ctx)

	return &types.PositionResponse{
		Position:      &position,
		LockedLpToken: totalLocked.AmountOf(ammtypes.GetPoolShareDenom(position.AmmPoolId)),
	}, nil
}

func (k Keeper) LiquidationPrice(goCtx context.Context, req *types.QueryLiquidationPriceRequest) (*types.QueryLiquidationPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	position, err := k.GetPosition(ctx, req.Address, req.PositionId)
	if err != nil {
		return nil, err
	}

	params := k.GetParams(ctx)

	debt := k.stableKeeper.GetDebt(ctx, position.GetPositionAddress())

	// Ensure position.LeveragedLpAmount is not zero to avoid division by zero
	if position.LeveragedLpAmount.IsZero() {
		return nil, types.ErrAmountTooLow
	}

	// lpTokenPrice * lpTokenAmount / totalDebt = params.SafetyFactor
	// lpTokenPrice = totalDebt * params.SafetyFactor / lpTokenAmount
	totalDebt := debt.GetTotalLiablities()
	baseCurrency, _ := k.assetProfileKeeper.GetUsdcDenom(ctx)
	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	liquidationPrice := params.SafetyFactor.MulInt(totalDebt).Mul(usdcDenomPrice).MulInt(ammtypes.OneShare).QuoInt(position.LeveragedLpAmount)

	return &types.QueryLiquidationPriceResponse{
		Price: liquidationPrice,
	}, nil
}
