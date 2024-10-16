package keeper

import (
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) OnPoolEnable(ctx sdk.Context, ammPool ammtypes.Pool) error {
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

func (k Keeper) OnPoolDisable(ctx sdk.Context, ammPool ammtypes.Pool) error {
	accountedPool, found := k.GetAccountedPool(ctx, ammPool.PoolId)
	if !found {
		return types.ErrPoolDoesNotExist
	}

	for _, nonAmmPoolToken := range accountedPool.NonAmmPoolTokens {
		if !nonAmmPoolToken.Amount.IsZero() {
			return fmt.Errorf("accounted pool have non-zero non amm pool balance left")
		}
	}

	k.RemoveAccountedPool(ctx, ammPool.PoolId)

	return nil
}

func (k Keeper) PerpetualUpdates(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool) error {
	// Get accounted pool
	accountedPool, found := k.GetAccountedPool(ctx, ammPool.PoolId)
	if !found {
		return types.ErrPoolDoesNotExist
	}

	// Accounted pool balance = amm pool + (long liability - long profit taking liability) - (long custody - long profit taking custody) + (short liability - short profit taking liability ) - ( short custody - short profit taking custody)
	// Accounted pool balance = amm pool + totalLiabilities - totalCustody + total profit taking custody - total profit taking liability
	for i, asset := range accountedPool.PoolAssets {
		ammBalance, err := ammPool.GetAmmPoolBalance(asset.Token.Denom)
		if err != nil {
			return err
		}
		totalLiabilities, totalCustody, totalTakeProfitCustody, totalTakeProfitLiabilities := perpetualPool.GetPerpetualPoolBalances(asset.Token.Denom)
		accountedPoolAmt := ammBalance.Add(totalLiabilities).Sub(totalCustody).Add(totalTakeProfitCustody).Sub(totalTakeProfitLiabilities)
		accountedPool.PoolAssets[i].Token = sdk.NewCoin(asset.Token.Denom, accountedPoolAmt)

		for j, nonAmmToken := range accountedPool.NonAmmPoolTokens {
			if nonAmmToken.Denom == asset.Token.Denom {
				accountedPool.NonAmmPoolTokens[j].Amount = accountedPoolAmt.Sub(ammBalance)
				break
			}
		}
	}

	// Update accounted pool
	k.SetAccountedPool(ctx, accountedPool)
	return nil
}

// Hooks wrapper struct for tvl keeper
type PerpetualHooks struct {
	k Keeper
}

var _ perpetualtypes.PerpetualHooks = PerpetualHooks{}

// Return the wrapper struct
func (k Keeper) PerpetualHooks() PerpetualHooks {
	return PerpetualHooks{k}
}

func (h PerpetualHooks) AfterEnablingPool(ctx sdk.Context, pool ammtypes.Pool) error {
	return h.k.OnPoolEnable(ctx, pool)
}

func (h PerpetualHooks) AfterDisablingPool(ctx sdk.Context, pool ammtypes.Pool) error {
	return h.k.OnPoolDisable(ctx, pool)
}

func (h PerpetualHooks) AfterPerpetualPositionOpen(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress) error {
	return h.k.PerpetualUpdates(ctx, ammPool, perpetualPool)
}

func (h PerpetualHooks) AfterPerpetualPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress) error {
	return h.k.PerpetualUpdates(ctx, ammPool, perpetualPool)
}

func (h PerpetualHooks) AfterPerpetualPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, perpetualPool perpetualtypes.Pool, sender sdk.AccAddress) error {
	return h.k.PerpetualUpdates(ctx, ammPool, perpetualPool)
}
