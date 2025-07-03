package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) CheckLowPoolHealthAndMinimumCustody(ctx sdk.Context, poolId uint64, openedPosition bool) error {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return errorsmod.Wrapf(types.ErrPoolDoesNotExist, "pool id %d", poolId)
	}

	params := k.GetParams(ctx)

	maxLiabilitiesRatioAllowed := math.LegacyZeroDec()
	if openedPosition {
		maxLiabilitiesRatioAllowed = params.PoolMaxLiabilitiesThreshold
	} else {
		maxLiabilitiesRatioAllowed = params.PoolMaxLiabilitiesThreshold.Add(params.ExitBuffer)
	}

	if !pool.BaseAssetLiabilitiesRatio.IsNil() && pool.BaseAssetLiabilitiesRatio.GTE(maxLiabilitiesRatioAllowed) {
		return errorsmod.Wrapf(types.ErrInvalidPosition, "pool (%d) base asset liabilities ratio (%s) too high for the operation", poolId, pool.BaseAssetLiabilitiesRatio.String())
	}
	if !pool.QuoteAssetLiabilitiesRatio.IsNil() && pool.QuoteAssetLiabilitiesRatio.GTE(maxLiabilitiesRatioAllowed) {
		return errorsmod.Wrapf(types.ErrInvalidPosition, "pool (%d) quote asset liabilities ratio (%s) too high for the operation", poolId, pool.QuoteAssetLiabilitiesRatio.String())
	}
	err := k.CheckMinimumCustodyAmt(ctx, poolId)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) CalculateLiabilitiesRatioByPosition(pool *types.Pool, ammPool ammtypes.Pool, position types.Position) math.LegacyDec {
	poolAssets := pool.GetPoolAssets(position)
	H := math.LegacyZeroDec()
	for _, asset := range *poolAssets {

		if asset.Liabilities.IsZero() {
			continue
		}

		ammBalance, err := ammPool.GetAmmPoolBalance(asset.AssetDenom)
		if err != nil {
			return math.LegacyZeroDec()
		}

		balance := ammBalance.ToLegacyDec()
		liabilities := asset.Liabilities.ToLegacyDec()

		if balance.Add(liabilities).IsZero() {
			return math.LegacyZeroDec()
		}
		H = liabilities.Quo(balance.Add(liabilities))
	}
	return H
}

func (k Keeper) UpdatePoolHealth(ctx sdk.Context, pool *types.Pool) error {
	ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
	if !found {
		return errors.New("amm pool not found while calculating pool health")
	}

	pool.BaseAssetLiabilitiesRatio = k.CalculateLiabilitiesRatioByPosition(pool, ammPool, types.Position_LONG)
	pool.QuoteAssetLiabilitiesRatio = k.CalculateLiabilitiesRatioByPosition(pool, ammPool, types.Position_SHORT)

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
		_, totalCustody := pool.GetPerpetualPoolBalances(ammPoolAsset.Token.Denom)
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

	totalLiabilities := math.ZeroInt()
	for _, poolAsset := range pool.PoolAssetsLong {
		// for long, liabilities will always be in base currency
		totalLiabilities = totalLiabilities.Add(poolAsset.Liabilities)
	}

	tradingAsset := ""
	for _, poolAsset := range pool.PoolAssetsLong {
		if poolAsset.AssetDenom != baseCurrency {
			tradingAsset = poolAsset.AssetDenom
			break
		}
	}

	_, tradingAssetPriceInBaseUnits, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, tradingAsset)
	if err != nil {
		return sdk.Coin{}, err
	}

	for _, poolAsset := range pool.PoolAssetsShort {
		// For short liabilities will be in trading asset
		baseCurrencyAmt := poolAsset.GetBigDecLiabilities().Mul(tradingAssetPriceInBaseUnits).Dec().TruncateInt()
		totalLiabilities = totalLiabilities.Add(baseCurrencyAmt)
	}
	return sdk.NewCoin(baseCurrency, totalLiabilities), nil
}
