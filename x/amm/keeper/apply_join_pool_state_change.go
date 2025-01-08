package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) ApplyJoinPoolStateChange(
	ctx sdk.Context,
	pool types.Pool,
	joiner sdk.AccAddress,
	numShares math.Int,
	joinCoins sdk.Coins,
	weightBalanceBonus elystypes.Dec34,
) error {
	if err := k.bankKeeper.SendCoins(ctx, joiner, sdk.MustAccAddressFromBech32(pool.GetAddress()), joinCoins); err != nil {
		return err
	}

	if err := k.MintPoolShareToAccount(ctx, pool, joiner, numShares); err != nil {
		return err
	}

	k.SetPool(ctx, pool)

	rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())

	if weightBalanceBonus.IsPositive() {
		// calculate treasury amounts to send as bonus
		weightBalanceBonusCoins := PortionCoins(joinCoins, weightBalanceBonus)
		for _, coin := range weightBalanceBonusCoins {
			treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, coin.Denom).Amount
			if treasuryTokenAmount.LT(coin.Amount) {
				// override coin amount by treasuryTokenAmount
				weightBalanceBonusCoins = weightBalanceBonusCoins.
					Sub(coin).                                        // remove the original coin
					Add(sdk.NewCoin(coin.Denom, treasuryTokenAmount)) // add the treasuryTokenAmount
			}
		}

		// send bonus tokens to recipient if positive
		if weightBalanceBonusCoins.IsAllPositive() {
			if err := k.bankKeeper.SendCoins(ctx, rebalanceTreasuryAddr, joiner, weightBalanceBonusCoins); err != nil {
				return err
			}
		}
	}

	types.EmitAddLiquidityEvent(ctx, joiner, pool.GetPoolId(), joinCoins)
	if k.hooks != nil {
		err := k.hooks.AfterJoinPool(ctx, joiner, pool, joinCoins, numShares)
		if err != nil {
			return err
		}
	}
	return nil
}
