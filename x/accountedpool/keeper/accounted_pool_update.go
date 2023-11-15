package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

// Get Amm Pool Balance
func (k Keeper) GetAmmPoolBalance(ammPool ammtypes.Pool, denom string) sdk.Int {
	for _, asset := range ammPool.PoolAssets {
		if asset.Token.Denom == denom {
			return asset.Token.Amount
		}
	}

	return sdk.ZeroInt()
}

func (k Keeper) GetMarginPoolBalancesByPosition(marginPool margintypes.Pool, denom string, position margintypes.Position) (sdk.Int, sdk.Int, sdk.Int) {
	poolAssets := marginPool.GetPoolAssets(position)

	for _, asset := range *poolAssets {
		if asset.AssetDenom == denom {
			return asset.AssetBalance, asset.Liabilities, asset.Custody
		}
	}

	return sdk.ZeroInt(), sdk.ZeroInt(), sdk.ZeroInt()
}

// Get Margin Pool Balance
func (k Keeper) GetMarginPoolBalances(marginPool margintypes.Pool, denom string) (sdk.Int, sdk.Int, sdk.Int) {
	assetBalanceLong, liabilitiesLong, custodyLong := k.GetMarginPoolBalancesByPosition(marginPool, denom, margintypes.Position_LONG)
	assetBalanceShort, liabilitiesShort, custodyShort := k.GetMarginPoolBalancesByPosition(marginPool, denom, margintypes.Position_SHORT)

	assetBalance := assetBalanceLong.Add(assetBalanceShort)
	liabilities := liabilitiesLong.Add(liabilitiesShort)
	custody := custodyLong.Add(custodyShort)

	return assetBalance, liabilities, custody
}

func (k Keeper) UpdateAccountedPool(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) error {
	poolId := ammPool.PoolId
	// Check if already exists
	exists := k.PoolExists(ctx, poolId)
	if !exists {
		return errors.New("pool doesn't exist!")
	}

	// Get accounted pool
	accountedPool, found := k.GetAccountedPool(ctx, poolId)
	if !found {
		return errors.New("pool doesn't exist!")
	}

	// Accounted Pool balance =
	// amm pool balance + margin pool balance + margin pool liabilties - margin pool custody
	// But not deducting custody amount here as the balance was already deducted through TakeCustody function.
	for i, asset := range accountedPool.PoolAssets {
		aBalance := k.GetAmmPoolBalance(ammPool, asset.Token.Denom)
		mBalance, mLiabiltiies, _ := k.GetMarginPoolBalances(marginPool, asset.Token.Denom)
		accountedAmt := aBalance.Add(mBalance).Add(mLiabiltiies)
		accountedPool.PoolAssets[i].Token = sdk.NewCoin(asset.Token.Denom, accountedAmt)
	}

	// Update accounted pool
	k.SetAccountedPool(ctx, accountedPool)
	return nil
}
