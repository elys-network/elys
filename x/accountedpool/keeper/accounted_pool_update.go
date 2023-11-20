package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

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
		aBalance, err := margintypes.GetAmmPoolBalance(ammPool, asset.Token.Denom)
		if err != nil {
			return err
		}
		mBalance, mLiabiltiies, _ := margintypes.GetMarginPoolBalances(marginPool, asset.Token.Denom)
		accountedAmt := aBalance.Add(mBalance).Add(mLiabiltiies)
		accountedPool.PoolAssets[i].Token = sdk.NewCoin(asset.Token.Denom, accountedAmt)
	}

	// Update accounted pool
	k.SetAccountedPool(ctx, accountedPool)
	return nil
}
