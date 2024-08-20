package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
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
	updatedLeveragePosition := types.QueryPosition{}

	pool, found := k.amm.GetPool(ctx, position.AmmPoolId)
	if !found {
		return nil, errorsmod.Wrap(ammtypes.ErrPoolNotFound, fmt.Sprintf("poolId: %d", position.AmmPoolId))
	}
	lp_price, err := pool.LpTokenPrice(ctx, k.oracleKeeper)
	if err != nil {
		return nil, err
	}

	lp_usd_price := position.LeveragedLpAmount.Mul(lp_price.TruncateInt())
	price := k.oracleKeeper.GetAssetPriceFromDenom(ctx, position.Collateral.Denom)
	updated_leverage := lp_usd_price.Quo(lp_usd_price.Sub(position.Liabilities.Mul(price.TruncateInt())))
	
	updatedLeveragePosition = types.QueryPosition{
		Position:        &position,
		UpdatedLeverage: updated_leverage,
	}

	commitments := k.commKeeper.GetCommitments(ctx, position.GetPositionAddress())
	totalLocked, _ := commitments.CommittedTokensLocked(ctx)

	return &types.PositionResponse{
		Position:      &updatedLeveragePosition,
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
	usdcDenomPrice := k.oracleKeeper.GetAssetPriceFromDenom(ctx, baseCurrency)
	liquidationPrice := params.SafetyFactor.MulInt(totalDebt).Mul(usdcDenomPrice).MulInt(ammtypes.OneShare).QuoInt(position.LeveragedLpAmount)

	return &types.QueryLiquidationPriceResponse{
		Price: liquidationPrice,
	}, nil
}
