package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckLowPoolHealth(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}

	if !pool.Health.IsNil() && pool.Health.LTE(k.GetPoolOpenThreshold(ctx)) {
		return errorsmod.Wrapf(types.ErrInvalidPosition, "pool (%d) health too low to open new positions", poolId)
	}
	return nil
}

func (k Keeper) CalculatePoolHealthByPosition(pool *types.Pool, ammPool ammtypes.Pool, position types.Position) sdk.Dec {
	poolAssets := pool.GetPoolAssets(position)
	H := sdk.NewDec(1)
	for _, asset := range *poolAssets {

		ammBalance, err := ammPool.GetAmmPoolBalance(asset.AssetDenom)
		if err != nil {
			return sdk.ZeroDec()
		}

		balance := ammBalance.ToLegacyDec()
		liabilities := asset.Liabilities.ToLegacyDec()

		if balance.Add(liabilities).IsZero() {
			return sdk.ZeroDec()
		}

		mul := balance.Quo(balance.Add(liabilities))
		H = H.Mul(mul)
	}
	return H
}

func (k Keeper) CalculatePoolHealth(ctx sdk.Context, pool *types.Pool) sdk.Dec {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return sdk.ZeroDec()
	}

	H := k.CalculatePoolHealthByPosition(pool, ammPool, types.Position_LONG)
	H = H.Mul(k.CalculatePoolHealthByPosition(pool, ammPool, types.Position_SHORT))

	return H
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) error {
	pool.Health = k.CalculatePoolHealth(ctx, pool)
	k.SetPool(ctx, *pool)

	return nil
}

// CheckMinimumCustodyAmt Should be called after opening positions and when real pool balance changes
func (k Keeper) CheckMinimumCustodyAmt(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}
	ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
	if err != nil {
		return err
	}
	for _, ammPoolAsset := range ammPool.PoolAssets {
		_, totalCustody, _, _ := pool.GetPerpetualPoolBalances(ammPoolAsset.Token.Denom)
		if ammPoolAsset.Token.Amount.LT(totalCustody) {
			return fmt.Errorf("real amm pool (id: %d) balance (%s) is less than total custody (%s)", poolId, ammPoolAsset.Token.String(), totalCustody.String())
		}
	}
	return nil
}
