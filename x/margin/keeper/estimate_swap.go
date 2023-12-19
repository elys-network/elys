package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

// Swap estimation using amm CalcOutAmtGivenIn function
func (k Keeper) EstimateSwap(ctx sdk.Context, tokenInAmount sdk.Coin, tokenOutDenom string, ammPool ammtypes.Pool) (sdk.Int, error) {
	marginEnabled := k.IsPoolEnabled(ctx, ammPool.PoolId)
	if !marginEnabled {
		return sdk.ZeroInt(), errorsmod.Wrap(types.ErrMarginDisabled, "Margin disabled pool")
	}

	tokensIn := sdk.Coins{tokenInAmount}
	// Estimate swap
	snapshot := k.amm.GetPoolSnapshotOrSet(ctx, ammPool)
	swapResult, err := k.amm.CalcOutAmtGivenIn(ctx, ammPool.PoolId, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, sdk.ZeroDec())
	if err != nil {
		return sdk.ZeroInt(), err
	}

	if swapResult.IsZero() {
		return sdk.ZeroInt(), types.ErrAmountTooLow
	}
	return swapResult.Amount, nil
}
