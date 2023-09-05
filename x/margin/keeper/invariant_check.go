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

	marginPool, found := k.GetPool(ctx, poolId)
	if !found {
		return errors.New("pool doesn't exist!")
	}

	address, err := sdk.AccAddressFromBech32(ammPool.GetAddress())
	if err != nil {
		return err
	}

	// bank balance should be ammPool balance + margin pool balance
	balances := k.bankKeeper.GetAllBalances(ctx, address)
	for _, balance := range balances {
		ammBalance, _ := k.GetAmmPoolBalance(ctx, ammPool, balance.Denom)
		marginBalance, _, _ := k.GetMarginPoolBalances(marginPool, balance.Denom)

		diff := ammBalance.Add(marginBalance).Sub(balance.Amount)
		if !diff.IsZero() {
			return errors.New("balance mismatch!")
		}
	}
	return nil
}

// Check if amm pool balance in bank module is correct
func (k Keeper) InvariantCheck(ctx sdk.Context) error {
	mtps := k.GetAllMTPs(ctx)
	for _, mtp := range mtps {
		ammPoolId := mtp.AmmPoolId
		err := k.AmmPoolBalanceCheck(ctx, ammPoolId)
		if err != nil {
			return err
		}
	}

	return nil
}
