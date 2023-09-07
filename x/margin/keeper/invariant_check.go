package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) AmmPoolBalanceCheck(ctx sdk.Context, poolId uint64) error {
	ammPool, found := k.amm.GetPool(ctx, poolId)
	if !found {
		return errors.New("pool doesn't exist!")
	}

	address, err := sdk.AccAddressFromBech32(ammPool.GetAddress())
	if err != nil {
		return err
	}

	mtpCollateralBalances := sdk.NewCoins()
	mtps := k.GetAllMTPs(ctx)
	for _, mtp := range mtps {
		ammPoolId := mtp.AmmPoolId
		if !k.OpenLongChecker.IsPoolEnabled(ctx, ammPoolId) {
			continue
		}

		if poolId != mtp.AmmPoolId {
			continue
		}

		mtpCollateralBalances = mtpCollateralBalances.Add(sdk.NewCoin(mtp.CollateralAsset, mtp.CollateralAmount))
	}

	// bank balance should be ammPool balance + collateral balance
	// TODO:
	// Need to think about correct algorithm of balance checking.
	// Important note.
	// AMM pool balance differs bank module balance
	balances := k.bankKeeper.GetAllBalances(ctx, address)
	for _, balance := range balances {
		ammBalance, _ := k.GetAmmPoolBalance(ctx, ammPool, balance.Denom)
		collateralAmt := mtpCollateralBalances.AmountOf(balance.Denom)

		diff := ammBalance.Add(collateralAmt).Sub(balance.Amount)
		if !diff.IsZero() {
			return errors.New("balance mismatch!")
		}
	}
	return nil
}

// Check if amm pool balance in bank module is correct
func (k Keeper) InvariantCheck(ctx sdk.Context) error {
	ammPools := k.amm.GetAllPool(ctx)
	for _, ammPool := range ammPools {
		err := k.AmmPoolBalanceCheck(ctx, ammPool.PoolId)
		if err != nil {
			return err
		}
	}

	return nil
}
