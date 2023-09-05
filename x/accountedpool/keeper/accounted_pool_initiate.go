package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) InitiateAccountedPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	poolId := ammPool.PoolId
	// Check if already exists
	exists := k.PoolExists(ctx, poolId)
	if exists {
		return errors.New("already existed pool!")
	}

	// Initiate pool
	accountedPool := types.AccountedPool{
		PoolId:      poolId,
		TotalShares: ammPool.TotalShares,
		PoolAssets:  []ammtypes.PoolAsset{},
		TotalWeight: ammPool.TotalWeight,
	}

	for _, asset := range ammPool.PoolAssets {
		accountedPool.PoolAssets = append(accountedPool.PoolAssets, asset)
	}
	// Set accounted pool
	k.SetAccountedPool(ctx, accountedPool)

	return nil
}
