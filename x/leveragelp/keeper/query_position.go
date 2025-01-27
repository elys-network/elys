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
	creator := sdk.MustAccAddressFromBech32(req.Address)
	position, err := k.GetPosition(ctx, creator, req.Id)
	if err != nil {
		return nil, err
	}
	updatedLeveragePosition, err := k.GetLeverageLpUpdatedLeverage(ctx, []*types.Position{&position})

	if err != nil {
		return nil, err
	}

	commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress())
	totalLocked, _ := commitments.CommittedTokensLocked(ctx)

	return &types.PositionResponse{
		Position:      updatedLeveragePosition[0],
		LockedLpToken: totalLocked.AmountOf(ammtypes.GetPoolShareDenom(position.AmmPoolId)),
	}, nil
}

func (k Keeper) LiquidationPrice(goCtx context.Context, req *types.QueryLiquidationPriceRequest) (*types.QueryLiquidationPriceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	creator := sdk.MustAccAddressFromBech32(req.Address)
	position, err := k.GetPosition(ctx, creator, req.PositionId)
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
	usdcDenomPrice, _ := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	liquidationPrice := usdcDenomPrice.MulLegacyDec(params.SafetyFactor).MulInt(totalDebt).MulInt(ammtypes.OneShare).QuoInt(position.LeveragedLpAmount)

	return &types.QueryLiquidationPriceResponse{
		Price: liquidationPrice.String(),
	}, nil
}
