package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// Swap estimation using amm CalcOutAmtGivenIn function
func (k Keeper) EstimateSwap(ctx sdk.Context, tokenInAmount sdk.Coin, tokenOutDenom string, ammPool ammtypes.Pool) (sdkmath.Int, error) {
	perpetualEnabled := k.IsPoolEnabled(ctx, ammPool.PoolId)
	if !perpetualEnabled {
		return sdkmath.ZeroInt(), errorsmod.Wrap(types.ErrPerpetualDisabled, "Perpetual disabled pool")
	}

	// Estimate swap
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, ammPool)
	tokensIn := sdk.Coins{tokenInAmount}
	swapResult, _, err := k.amm.CalcOutAmtGivenIn(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, sdkmath.LegacyZeroDec())
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	if swapResult.IsZero() {
		return sdkmath.ZeroInt(), types.ErrAmountTooLow
	}
	return swapResult.Amount, nil
}
