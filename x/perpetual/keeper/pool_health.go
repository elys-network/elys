package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckLowPoolHealthAndMinimumCustody(ctx sdk.Context, poolId uint64) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}

	minimumThreshold := k.GetPoolOpenThreshold(ctx)
	if !pool.Health.IsNil() && pool.Health.LTE(minimumThreshold) {
		return errorsmod.Wrapf(types.ErrInvalidPosition, "pool (%d) health too low to open new positions", poolId)
	}
	err := k.CheckMinimumCustodyAmt(ctx, poolId)
	if err != nil {
		return err
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

func (k Keeper) GetPoolTotalBaseCurrencyLiabilities(ctx sdk.Context, pool types.Pool) (sdk.Coin, error) {
	// retrieve base currency denom
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return sdk.Coin{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	totalLiabilities := math.LegacyZeroDec()
	for _, poolAsset := range pool.PoolAssetsLong {
		// for long, liabilities will always be in base currency
		totalLiabilities = totalLiabilities.Add(poolAsset.Liabilities.ToLegacyDec())
	}

	tradingAsset := ""
	for _, poolAsset := range pool.PoolAssetsLong {
		if poolAsset.AssetDenom != baseCurrency {
			tradingAsset = poolAsset.AssetDenom
			break
		}
	}

	tradingAssetPrice, err := k.GetAssetPrice(ctx, tradingAsset)
	if err != nil {
		return sdk.Coin{}, err
	}

	for _, poolAsset := range pool.PoolAssetsShort {
		// For short liabilities will be in trading asset
		baseCurrencyAmt := poolAsset.Liabilities.ToLegacyDec().Mul(tradingAssetPrice)
		totalLiabilities = totalLiabilities.Add(baseCurrencyAmt)
	}
	return sdk.NewCoin(baseCurrency, totalLiabilities.TruncateInt()), nil
}
