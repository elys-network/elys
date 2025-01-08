package keeper

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapOutAmtGivenIn is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
func (k Keeper) SwapOutAmtGivenIn(
	ctx sdk.Context, poolId uint64,
	oracleKeeper types.OracleKeeper,
	snapshot *types.Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdkmath.LegacyDec,
	weightBreakingFeePerpetualFactor sdkmath.LegacyDec,
) (tokenOut sdk.Coin, slippage, slippageAmount, weightBalanceBonus, oracleOutAmount elystypes.Dec34, swapFeeFinal sdkmath.LegacyDec, err error) {
	ammPool, found := k.GetPool(ctx, poolId)
	if !found {
		return sdk.Coin{}, elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), elystypes.ZeroDec34(), sdkmath.LegacyZeroDec(), fmt.Errorf("invalid pool: %d", poolId)
	}
	params := k.GetParams(ctx)
	return ammPool.SwapOutAmtGivenIn(ctx, oracleKeeper, snapshot, tokensIn, tokenOutDenom, swapFee, k.accountedPoolKeeper, weightBreakingFeePerpetualFactor, params)
}
