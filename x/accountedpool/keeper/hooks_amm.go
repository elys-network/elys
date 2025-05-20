package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
)

func (k Keeper) UpdateAccountedPoolOnAmmChange(ctx sdk.Context, ammPool ammtypes.Pool) error {
	poolId := ammPool.PoolId
	accountedPool, found := k.GetAccountedPool(ctx, poolId)
	if !found {
		// possible that pool does not exist on accounted pool
		return nil
	}

	for _, ammPoolAsset := range ammPool.PoolAssets {

		nonAmmTokenBalance := math.ZeroInt()

		for _, nonAmmPoolToken := range accountedPool.NonAmmPoolTokens {
			if nonAmmPoolToken.Denom == ammPoolAsset.Token.Denom {
				nonAmmTokenBalance = nonAmmPoolToken.Amount
				break
			}
		}

		for j, accountedPoolAsset := range accountedPool.TotalTokens {

			if ammPoolAsset.Token.Denom == accountedPoolAsset.Denom {
				updatedAccountedPoolAsset := ammPoolAsset
				updatedAccountedPoolAsset.Token.Amount = updatedAccountedPoolAsset.Token.Amount.Add(nonAmmTokenBalance)
				accountedPool.TotalTokens[j] = updatedAccountedPoolAsset.Token
				break
			}

		}
	}
	// Set accounted pool
	k.SetAccountedPool(ctx, accountedPool)

	//ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventAmmChanges,
	//	sdk.NewAttribute("pool_id", strconv.FormatUint(poolId, 10)),
	//	sdk.NewAttribute("non_amm_token_balance", sdk.Coins(accountedPool.NonAmmPoolTokens).String()),
	//	sdk.NewAttribute("total_tokens", sdk.Coins(accountedPool.TotalTokens).String()),
	//))

	return nil
}

type AmmHooks struct {
	k Keeper
}

var _ ammtypes.AmmHooks = AmmHooks{}

func (k Keeper) AmmHooks() AmmHooks {
	return AmmHooks{k}
}

func (h AmmHooks) AfterPoolCreated(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool) error {
	return nil
}

func (h AmmHooks) AfterJoinPool(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, enterCoins sdk.Coins, shareOutAmount math.Int) error {
	return h.k.UpdateAccountedPoolOnAmmChange(ctx, pool)
}

func (h AmmHooks) AfterExitPool(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, shareInAmount math.Int, exitCoins sdk.Coins) error {
	return h.k.UpdateAccountedPoolOnAmmChange(ctx, pool)
}

func (h AmmHooks) AfterSwap(ctx sdk.Context, sender sdk.AccAddress, pool ammtypes.Pool, input sdk.Coins, output sdk.Coins) error {
	return h.k.UpdateAccountedPoolOnAmmChange(ctx, pool)
}
