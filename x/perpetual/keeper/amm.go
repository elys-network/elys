package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammkeeper "github.com/elys-network/elys/v6/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) GetAmmPool(ctx sdk.Context, poolId uint64) (ammtypes.Pool, error) {
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return ammPool, errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}
	return ammPool, nil
}

func getWeightBreakingFee(weightBalanceBonus osmomath.BigDec) osmomath.BigDec {
	// when weightBalanceBonus is 0, then breaking fee is also 0
	// when it's > 0, then breaking fee is still 0
	// when it's < 0, breaking fee is it's negative
	if weightBalanceBonus.IsNegative() {
		return weightBalanceBonus.Neg()
	} else {
		return osmomath.ZeroBigDec()
	}
}

// Swap estimation using amm CalcOutAmtGivenIn function
func (k Keeper) EstimateSwapGivenIn(ctx sdk.Context, tokenInAmount sdk.Coin, tokenOutDenom string, ammPool ammtypes.Pool, owner string) (math.Int, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, math.LegacyDec, math.LegacyDec, error) {
	if tokenInAmount.IsZero() {
		return math.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, osmomath.BigDec{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenInAmount is zero for EstimateSwapGivenIn")
	}
	params := k.GetParams(ctx)

	// Estimate swap
	snapshot := k.amm.GetPoolWithAccountedBalance(ctx, ammPool.PoolId)
	tokensIn := sdk.Coins{tokenInAmount}
	tokenOut, slippage, slippageAmount, weightBalanceBonus, _, _, err := k.amm.SwapOutAmtGivenIn(ctx, ammPool.PoolId, k.oracleKeeper, snapshot, tokensIn, tokenOutDenom, osmomath.ZeroBigDec(), params.GetBigDecWeightBreakingFeeFactor(), osmomath.ZeroBigDec())
	if err != nil {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}

	if tokenOut.IsZero() {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(types.ErrAmountTooLow, "tokenOut is zero for swap (EstimateSwapGivenIn) for in %s and out denom %s", tokenInAmount.String(), tokenOutDenom)
	}

	// send for query purpose only
	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)
	perpetualFees := ammtypes.ApplyDiscount(params.GetBigDecPerpetualSwapFee(), tier.GetBigDecDiscount())
	takersFee := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()

	return tokenOut.Amount, slippage, slippageAmount, getWeightBreakingFee(weightBalanceBonus), perpetualFees.Dec(), takersFee.Dec(), nil
}

// Swap estimation using amm CalcInAmtGivenOut function
func (k Keeper) EstimateSwapGivenOut(ctx sdk.Context, tokenOutAmount sdk.Coin, tokenInDenom string, ammPool ammtypes.Pool, owner string) (math.Int, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, osmomath.BigDec, math.LegacyDec, math.LegacyDec, error) {
	if tokenOutAmount.IsZero() {
		return math.Int{}, osmomath.BigDec{}, osmomath.BigDec{}, osmomath.BigDec{}, osmomath.BigDec{}, math.LegacyDec{}, math.LegacyDec{}, fmt.Errorf("tokenOutAmount is zero for EstimateSwapGivenOut")
	}
	params := k.GetParams(ctx)
	tokensOut := sdk.Coins{tokenOutAmount}

	// Estimate swap
	snapshot := k.amm.GetPoolWithAccountedBalance(ctx, ammPool.PoolId)
	tokenIn, slippage, slippageAmount, weightBalanceBonus, oracleIn, _, err := k.amm.SwapInAmtGivenOut(ctx, ammPool.PoolId, k.oracleKeeper, snapshot, tokensOut, tokenInDenom, osmomath.ZeroBigDec(), params.GetBigDecWeightBreakingFeeFactor(), osmomath.ZeroBigDec())
	if err != nil {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(err, "unable to swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}

	if tokenIn.IsZero() {
		return math.ZeroInt(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), math.LegacyDec{}, math.LegacyDec{}, errorsmod.Wrapf(types.ErrAmountTooLow, "tokenIn is zero for swap (EstimateSwapGivenOut) for out %s and in denom %s", tokenOutAmount.String(), tokenInDenom)
	}

	// send for query purpose only
	addr, err := sdk.AccAddressFromBech32(owner)
	if err != nil {
		addr = sdk.AccAddress{}
	}
	_, tier := k.tierKeeper.GetMembershipTier(ctx, addr)
	perpetualFees := ammtypes.ApplyDiscount(params.GetBigDecPerpetualSwapFee(), tier.GetBigDecDiscount())
	takersFee := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees()

	return tokenIn.Amount, slippage, slippageAmount, getWeightBreakingFee(weightBalanceBonus), oracleIn, perpetualFees.Dec(), takersFee.Dec(), nil
}

// CalculatePerpetualFees calculates the perpetual fees, slippage fees, weight breaking fees, and taker fees for a swap.
// Pass calculatePerpAndTakerFees as false to exclude perp and taker fees for the swap
func (k Keeper) CalculatePerpetualFees(
	ctx sdk.Context,
	poolIsOracle bool,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
	slippageAmount osmomath.BigDec,
	weightBreakingFee osmomath.BigDec,
	perpetualFees math.LegacyDec,
	takersFee math.LegacyDec,
	oracleInAmount osmomath.BigDec,
	isSwapGivenIn bool,
	calculatePerpAndTakerFees bool,
) (perpFees types.PerpetualFees) {

	perpFeesCoins := sdk.Coins{}
	takerFeesCoins := sdk.Coins{}

	if calculatePerpAndTakerFees {

		// Determine the source of fees based on isSwapGivenIn
		takeFeesFrom := sdk.Coins{tokenIn}
		if !isSwapGivenIn && poolIsOracle {
			takeFeesFrom = sdk.NewCoins(sdk.NewCoin(tokenIn.Denom, oracleInAmount.Dec().TruncateInt()))
		}

		// Calculate perpetual fees in USD
		if perpetualFees.IsPositive() {
			perpFeesCoins = ammkeeper.PortionCoins(takeFeesFrom, osmomath.BigDecFromDec(perpetualFees))
		}

		// Calculate taker fees in USD
		if takersFee.IsPositive() {
			takerFeesCoins = ammkeeper.PortionCoins(takeFeesFrom, osmomath.BigDecFromDec(takersFee))
		}
	}

	// Calculate slippage amount in USD
	slippageCoins := sdk.Coins{}
	if isSwapGivenIn {
		slippageCoins = append(slippageCoins, sdk.NewCoin(tokenOut.Denom, slippageAmount.Dec().RoundInt()))
	} else {
		slippageCoins = append(slippageCoins, sdk.NewCoin(tokenIn.Denom, slippageAmount.Dec().RoundInt()))
	}

	// Calculate weight breaking fees in USD
	weightBreakingFeesCoins := sdk.Coins{}
	if !weightBreakingFee.IsZero() {
		var weightBreakingFeeAmount math.Int
		if isSwapGivenIn {
			weightBreakingFeeAmount = osmomath.BigDecFromSDKInt(tokenIn.Amount).Mul(weightBreakingFee).Dec().RoundInt()
		} else {
			weightBreakingFeeAmount = oracleInAmount.Mul(weightBreakingFee).Dec().RoundInt()
		}
		weightBreakingFeesCoins = append(weightBreakingFeesCoins, sdk.NewCoin(tokenIn.Denom, weightBreakingFeeAmount))
	}

	perpFees = types.NewPerpetualFees(
		perpFeesCoins,
		slippageCoins,
		weightBreakingFeesCoins,
		takerFeesCoins,
	)

	return perpFees
}
