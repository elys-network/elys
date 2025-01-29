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
		TotalTokens:      sdk.NewCoins(),
		NonAmmPoolTokens: sdk.NewCoins(),
	}

	nonAmmPoolTokens := make([]sdk.Coin, len(ammPool.PoolAssets))

	for i, poolAsset := range ammPool.PoolAssets {
		accountedPool.TotalTokens = append(accountedPool.TotalTokens, poolAsset.Token)
		nonAmmPoolTokens[i] = sdk.NewCoin(poolAsset.Token.Denom, math.ZeroInt())
	}
	accountedPool.NonAmmPoolTokens = nonAmmPoolTokens
	// Set accounted pool
	k.SetAccountedPool(ctx, accountedPool)

	//ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventLeverageLpEnable,
	//	sdk.NewAttribute("pool_id", strconv.FormatUint(poolId, 10)),
	//	sdk.NewAttribute("initial_tokens", sdk.Coins(accountedPool.TotalTokens).String()),
	//))

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

	//ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventLeverageLpDisable,
	//	sdk.NewAttribute("pool_id", strconv.FormatUint(ammPool.PoolId, 10)),
	//))

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

func (h LeverageLpHooks) AfterLeverageLpPositionOpen(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	return nil
}

func (h LeverageLpHooks) AfterLeverageLpPositionClose(ctx sdk.Context, _ sdk.AccAddress, ammPool ammtypes.Pool) error {
	return h.k.UpdateAccountedPoolOnAmmChange(ctx, ammPool)
}

func (h LeverageLpHooks) AfterLeverageLpPositionOpenConsolidate(ctx sdk.Context, sender sdk.AccAddress, ammPool ammtypes.Pool) error {
	return nil
}
