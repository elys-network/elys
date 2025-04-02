package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) JoinPoolEstimation(goCtx context.Context, req *types.QueryJoinPoolEstimationRequest) (*types.QueryJoinPoolEstimationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	tokensIn, sharesOut, slippage, weightBalanceBonus, swapFee, takerFees, weightRewardAmount, err := k.JoinPoolEst(ctx, req.PoolId, req.AmountsIn)
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetPoolShareDenom(req.PoolId)
	return &types.QueryJoinPoolEstimationResponse{
		ShareAmountOut:            sdk.NewCoin(shareDenom, sharesOut),
		AmountsIn:                 tokensIn,
		Slippage:                  slippage.String(),
		WeightBalanceRatio:        weightBalanceBonus.String(),
		SwapFee:                   swapFee.String(),
		TakerFee:                  takerFees.String(),
		WeightBalanceRewardAmount: weightRewardAmount,
	}, nil
}

func (k Keeper) JoinPoolEst(
	ctx sdk.Context,
	poolId uint64,
	tokenInMaxs sdk.Coins,
) (tokensIn sdk.Coins, sharesOut math.Int, slippage elystypes.Dec34, weightBalanceBonus elystypes.Dec34, swapFee elystypes.Dec34, takerFeesFinal elystypes.Dec34, weightRewardAmount sdk.Coin, err error) {
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, math.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			numShares, tokensIn, err := pool.CalcJoinPoolNoSwapShares(tokenInMaxs)
			if err != nil {
				return tokensIn, numShares, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, err
			}
		}

		params := k.GetParams(ctx)
		takerFees := elystypes.NewDec34FromLegacyDec(k.parameterKeeper.GetParams(ctx).TakerFees)
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		cacheCtx, _ := ctx.CacheContext()
		tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokensIn, params, takerFees)
		if err != nil {
			return nil, math.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, err
		}

		return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, sdk.Coin{}, nil
	}

	params := k.GetParams(ctx)
	takerFees := elystypes.NewDec34FromLegacyDec(k.parameterKeeper.GetParams(ctx).TakerFees)
	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	cacheCtx, _ := ctx.CacheContext()
	tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, _, err := pool.JoinPool(cacheCtx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs, params, takerFees)
	if err != nil {
		return nil, math.ZeroInt(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdk.Coin{}, err
	}

	var otherAsset types.PoolAsset
	bonusTokenAmount := math.ZeroInt()
	// Check treasury and update weightBalance
	if weightBalanceBonus.IsPositive() && tokensJoined.Len() == 1 {
		rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		for _, asset := range pool.PoolAssets {
			if asset.Token.Denom == tokensJoined[0].Denom {
				continue
			}
			otherAsset = asset
		}
		treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, otherAsset.Token.Denom).Amount
		bonusTokenAmount = weightBalanceBonus.MulInt(tokensJoined[0].Amount).ToInt()

		if treasuryTokenAmount.LT(bonusTokenAmount) {
			bonusTokenAmount = treasuryTokenAmount
		}
	}
	rewards := sdk.Coin{}
	if otherAsset.Token.Denom != "" && bonusTokenAmount.IsPositive() {
		rewards = sdk.NewCoin(otherAsset.Token.Denom, bonusTokenAmount)
	}

	return tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, takerFeesFinal, rewards, nil
}
