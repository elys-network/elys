package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) UpdateAccountedPool(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) error {
	poolId := ammPool.PoolId
	// Check if already exists
	exists := k.PoolExists(ctx, poolId)
	if !exists {
		return types.ErrPoolDoesNotExist
	}

	// Get accounted pool
	accountedPool, found := k.GetAccountedPool(ctx, poolId)
	if !found {
		return types.ErrPoolDoesNotExist
	}

	// Accounted Pool balance =
	// amm pool balance + perpetual pool balance + perpetual pool liabilties - perpetual pool custody
	// But not deducting custody amount here as the balance was already deducted through TakeCustody function.
	for i, asset := range accountedPool.PoolAssets {
		aBalance, err := perpetualtypes.GetAmmPoolBalance(ammPool, asset.Token.Denom)
		if err != nil {
			return err
		}
		mBalance, mLiabiltiies, _ := perpetualtypes.GetPerpetualPoolBalances(perpetualPool, asset.Token.Denom)
		accountedAmt := aBalance.Add(mBalance).Add(mLiabiltiies)
		accountedPool.PoolAssets[i].Token = sdk.NewCoin(asset.Token.Denom, accountedAmt)
	}

	// Update accounted pool
	k.SetAccountedPool(ctx, accountedPool)
	return nil
}
