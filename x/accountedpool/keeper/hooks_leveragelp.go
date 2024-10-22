package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	leveragelptypes "github.com/elys-network/elys/x/leveragelp/types"
)

func (k Keeper) OnLeverageLpPoolEnable(ctx sdk.Context, ammPool ammtypes.Pool) error {
	poolId := ammPool.PoolId
	// Check if already exists
	exists := k.PoolExists(ctx, poolId)
	if exists {
		return types.ErrPoolAlreadyExist
	}

	// Initiate pool
	accountedPool := types.AccountedPool{
		PoolId:           poolId,
		TotalShares:      ammPool.TotalShares,
		PoolAssets:       []ammtypes.PoolAsset{},
		TotalWeight:      ammPool.TotalWeight,
		NonAmmPoolTokens: sdk.NewCoins(),
	}

	nonAmmPoolTokens := make([]sdk.Coin, len(ammPool.PoolAssets))

	for i, poolAsset := range ammPool.PoolAssets {
		accountedPool.PoolAssets = append(accountedPool.PoolAssets, poolAsset)
		nonAmmPoolTokens[i] = sdk.NewCoin(poolAsset.Token.Denom, math.ZeroInt())
	}
	accountedPool.NonAmmPoolTokens = nonAmmPoolTokens
	// Set accounted pool
	k.SetAccountedPool(ctx, accountedPool)

	return nil
}

func (k Keeper) OnLeverageLpPoolDisable(ctx sdk.Context, ammPool ammtypes.Pool) error {
	accountedPool, found := k.GetAccountedPool(ctx, ammPool.PoolId)
	if !found {
		return types.ErrPoolDoesNotExist
	}

	// these are the balances updated on perpetual position changes
	for _, nonAmmPoolToken := range accountedPool.NonAmmPoolTokens {
		if !nonAmmPoolToken.Amount.IsZero() {
			return fmt.Errorf("all positions are not closed; accounted pool have non-zero non amm pool balance left")
		}
	}

	k.RemoveAccountedPool(ctx, ammPool.PoolId)

	return nil
}

type LeverageLpHooks struct {
	k Keeper
}

var _ leveragelptypes.LeverageLpHooks = LeverageLpHooks{}

// Return the wrapper struct
func (k Keeper) LeverageLpHooks() LeverageLpHooks {
	return LeverageLpHooks{k}
}

func (h LeverageLpHooks) AfterEnablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	return h.k.OnLeverageLpPoolEnable(ctx, ammPool)
}

func (h LeverageLpHooks) AfterDisablingPool(ctx sdk.Context, ammPool ammtypes.Pool) error {
	return h.k.OnLeverageLpPoolDisable(ctx, ammPool)
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress) error {
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, sender sdk.AccAddress) error {
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress) error {
	return nil
}
