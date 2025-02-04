package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) ExitPoolEstimation(goCtx context.Context, req *types.QueryExitPoolEstimationRequest) (*types.QueryExitPoolEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	exitCoins, weightBalanceBonus, slippage, err := k.ExitPoolEst(ctx, req.PoolId, req.ShareAmountIn, req.TokenOutDenom)
	if err != nil {
		return nil, err
	}

	return &types.QueryExitPoolEstimationResponse{
		AmountsOut:         exitCoins,
		WeightBalanceRatio: weightBalanceBonus,
		Slippage:           slippage,
	}, nil
}

func (k Keeper) ExitPoolEst(
	ctx sdk.Context,
	poolId uint64,
	shareInAmount math.Int,
	tokenOutDenom string,
) (exitCoins sdk.Coins, weightBalanceBonus math.LegacyDec, slippage math.LegacyDec, err error) {
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), types.ErrInvalidPoolId
	}

	totalSharesAmount := pool.GetTotalShares()
	if shareInAmount.GTE(totalSharesAmount.Amount) {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit >= the number of shares contained in the pool.")
	} else if shareInAmount.LTE(math.ZeroInt()) {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), errorsmod.Wrapf(types.ErrInvalidMathApprox, "Trying to exit a negative amount of shares")
	}

	params := k.GetParams(ctx)
	exitCoins, weightBalanceBonus, slippage, err = pool.CalcExitPoolCoinsFromShares(ctx, k.oracleKeeper, k.accountedPoolKeeper, shareInAmount, tokenOutDenom, params, true)
	if err != nil {
		return sdk.Coins{}, math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	return exitCoins, weightBalanceBonus, slippage, nil
}
